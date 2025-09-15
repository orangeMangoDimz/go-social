package httpserver

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/orangeMangoDimz/go-social/docs"
	"github.com/orangeMangoDimz/go-social/internal/auth"
	"github.com/orangeMangoDimz/go-social/internal/config"
	"github.com/orangeMangoDimz/go-social/internal/db"
	"github.com/orangeMangoDimz/go-social/internal/env"
	"github.com/orangeMangoDimz/go-social/internal/mailer"
	"github.com/orangeMangoDimz/go-social/internal/ratelimiter"
	authHandler "github.com/orangeMangoDimz/go-social/internal/server/http/handler/auth"
	healthHandler "github.com/orangeMangoDimz/go-social/internal/server/http/handler/health"
	postsHandler "github.com/orangeMangoDimz/go-social/internal/server/http/handler/posts"
	usersHandler "github.com/orangeMangoDimz/go-social/internal/server/http/handler/users"
	"github.com/orangeMangoDimz/go-social/internal/storage/cache"
	"github.com/orangeMangoDimz/go-social/internal/storage/postgres"
	"github.com/redis/go-redis/v9"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"
)

// Configuration loaders
func loadConfig() config.Config {
	return config.Config{
		Addr:        env.GetString("ADDR", ":8000"),
		Db:          loadDbConfig(),
		Env:         env.GetString("ENV", "development"),
		ApiURL:      env.GetString("EXTERNAL_URL", "localhost:8000"),
		FrontendURL: env.GetString("FRONTEND_URL", "http://localhost:3000"),
		Mail:        loadMailConfig(),
		Auth:        loadAuthConfig(),
		RedisCfg:    loadRedisConfig(),
		RateLimiter: loadRateLimiterConfig(),
	}
}

func loadDbConfig() config.DbConfig {
	return config.DbConfig{
		Addr:         env.GetString("DB_ADDR", "postgres://postgres:root@localhost/social?sslmode=disable"),
		MaxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		MaxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		MaxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	}
}

func loadMailConfig() config.MailConfig {
	return config.MailConfig{
		Exp:       time.Hour * 24 * 3, // 3 days
		FromEmail: env.GetString("FROM_EMAIL", ""),
		SendGrid: config.SendGridConfig{
			ApiKey: env.GetString("SENDGRID_API_KEY", ""),
		},
		MailTrap: config.MailTrapConfig{
			ApiKey: env.GetString("MAILTRAP_API_KEY", ""),
		},
	}
}

func loadAuthConfig() config.AuthConfig {
	return config.AuthConfig{
		Basic: config.BasicConfig{
			User: env.GetString("AUTH_BASIC_USER", "admin"),
			Pass: env.GetString("AUTH_BASIC_PASS", "admin"),
		},
		Token: config.TokenConfig{
			Secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
			Exp:    time.Hour * 24 * 3, // 3 days
			Iss:    "gophersocial",
		},
	}
}

func loadRedisConfig() config.RedisConfig {
	return config.RedisConfig{
		Addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
		Password: env.GetString("REDIS_PASSWORD", ""),
		Db:       env.GetInt("REDIS_DB", 0),
		Enabled:  env.GetBool("REDIS_ENABLED", false),
	}
}

func loadRateLimiterConfig() ratelimiter.Config {
	return ratelimiter.Config{
		RequestPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20),
		TimeFrame:           time.Second * 5,
		Enabled:             env.GetBool("RATELIMITER_ENABLED", true),
	}
}

// Component initializers
func initDatabase(cfg config.DbConfig, logger *zap.SugaredLogger) *sql.DB {
	db, err := db.New(cfg.Addr, cfg.MaxOpenConns, cfg.MaxIdleConns, cfg.MaxIdleTime)
	if err != nil {
		logger.Fatal("Failed to initialize database:", err)
	}
	logger.Info("Database connection pool established")
	return db
}

func initCache(cfg config.RedisConfig, logger *zap.SugaredLogger) *redis.Client {
	if !cfg.Enabled {
		return nil
	}
	rdb := cache.NewRedisClient(cfg.Addr, cfg.Password, cfg.Db)
	logger.Info("Redis connection established")
	return rdb
}

func initMailer(cfg config.MailConfig, logger *zap.SugaredLogger) mailer.Client {
	client, err := mailer.NewMailTrapClient(cfg.MailTrap.ApiKey, cfg.FromEmail)
	if err != nil {
		logger.Fatal("Failed to initialize mailer:", err)
	}
	return client
}

// NewApp creates and configures a new Application instance
func NewApp() (*sql.DB, *Application) {
	// Initialize logger first
	logger := zap.Must(zap.NewProduction()).Sugar()

	// Load configuration
	config := loadConfig()

	// Initialize components
	database := initDatabase(config.Db, logger)
	cacheClient := initCache(config.RedisCfg, logger)
	mailClient := initMailer(config.Mail, logger)

	// Initialize auth and rate limiter
	jwtAuth := auth.NewJWTAuthenticator(
		config.Auth.Token.Secret,
		config.Auth.Token.Iss,
		config.Auth.Token.Iss,
	)

	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		config.RateLimiter.RequestPerTimeFrame,
		config.RateLimiter.TimeFrame,
	)

	// Build application
	app := Application{
		Config:        config,
		Store:         postgres.NewStore(database),
		CacheStorage:  cache.NewRedisStorage(cacheClient),
		Logger:        logger,
		Mail:          mailClient,
		Authenticator: jwtAuth,
		RateLimiter:   rateLimiter,
	}

	return database, &app
}

func (app *Application) Mount(version string) http.Handler {
	r := chi.NewRouter()

	// A good base middleware stack
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{env.GetString("CORS_ALLOWED_ORIGIN", "")},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	if app.Config.RateLimiter.Enabled {
		r.Use(app.RateLimiterMiddleware)
	}

	// Set a timeout value on the request context (ctx), that will signal
	// through ctx.Done() that the request has timed out and further
	// processing should be stopped.
	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Group(healthHandler.RegisterRoute(app, app.Config, version))

		docsURL := fmt.Sprintf("%s/swagger/doc.json", app.Config.Addr)
		r.Get("/swagger/*", httpSwagger.Handler(httpSwagger.URL(docsURL)))
		r.Route("/posts", postsHandler.RegisterRoute(app, app.Services.PostService, app.Services.CommentService, *app.Logger))
		r.Route("/users", usersHandler.RegisterRoute(app, app.Services.UsersService, app.Services.FollowerService, *app.Logger))
		// Authentication routes
		r.Route("/authentication", authHandler.RegisterRoute(app.Services.UsersService, *app.Logger, app.Mail, app.Config, app.Authenticator))
	})

	return r
}

func (app *Application) Run(handler http.Handler, version string) error {
	// Docs
	docs.SwaggerInfo.Version = version
	docs.SwaggerInfo.Host = app.Config.ApiURL
	docs.SwaggerInfo.BasePath = "/v1"

	server := http.Server{
		Addr:         app.Config.Addr,
		Handler:      handler,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	// Channel to handle graceful shutdown errors
	shutdown := make(chan error)

	// Goroutine for graceful shutdown handling
	// Listens for OS signals (SIGINT/SIGTERM) and initiates server shutdown
	go func() {
		quit := make(chan os.Signal, 1)

		// Register the channel to receive specific OS signals
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		s := <-quit // Block until a signal is received

		// Create context with 5-second timeout for graceful shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.Logger.Infow("signal caught", "signal", s.String())

		// Send shutdown result to main goroutine
		shutdown <- server.Shutdown(ctx)
	}()

	app.Logger.Infow("server has started", "addr", app.Config.Addr, "env", app.Config.Env)

	// Start the HTTP server (blocking call)
	err := server.ListenAndServe()
	if !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	// Wait for shutdown completion and handle any shutdown errors
	err = <-shutdown
	if err != nil {
		return err
	}

	app.Logger.Info("Server stopped gracefully")
	return nil
}
