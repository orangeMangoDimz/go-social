package main

import (
	"time"

	"github.com/orangeMangoDimz/go-social/internal/auth"
	"github.com/orangeMangoDimz/go-social/internal/db"
	"github.com/orangeMangoDimz/go-social/internal/env"
	"github.com/orangeMangoDimz/go-social/internal/mailer"
	"github.com/orangeMangoDimz/go-social/internal/ratelimiter"
	"github.com/orangeMangoDimz/go-social/store"
	"github.com/orangeMangoDimz/go-social/store/cache"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const VERSION = ""

//	@title			Gopher Social API
//	@description	API for GopherSocial, a social network for gophers
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath					/v1
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Enter your bearer token in the format **Bearer &lt;token&gt;**

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						X-API-KEY
//	@description				API Key for authorization

//	@description

func main() {
	cfg := dbConfig{
		addr:         env.GetString("DB_ADDR", "postgres://postgres:root@localhost/social?sslmode=disable"),
		maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	}

	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	// Database
	db, err := db.New(cfg.addr, cfg.maxOpenConns, cfg.maxIdleConns, cfg.maxIdleTime)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("DB Database connection pool established")
	mail_config := mailConfig{
		exp:       time.Hour * 24 * 3, // 3 days
		fromEmail: env.GetString("FROM_EMAIL", ""),
		sendGrid: SendGridConfig{
			apiKey: env.GetString("SENDGRID_API_KEY", ""),
		},
		mailTrap: mailTrapConfig{
			apiKey: env.GetString("MAILTRAP_API_KEY", ""),
		},
	}

	// mailer := mailer.NewSendGrid(mail_config.sendGrid.apiKey, mail_config.fromEmail)
	mailtrap, err := mailer.NewMailTrapClient(mail_config.mailTrap.apiKey, mail_config.fromEmail)
	if err != nil {
		logger.Fatal(err)
	}

	rateLimiterCfg := ratelimiter.Config{
		RequestPerTimeFrame: env.GetInt("RATELIMITER_REQUESTS_COUNT", 20),
		TimeFrame:           time.Second * 5,
		Enabled:             env.GetBool("RATELIMITER_ENABLED", true),
	}

	config := config{
		addr:        env.GetString("ADDR", ":8000"),
		db:          cfg,
		env:         env.GetString("ENV", "development"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8000"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:3000"),
		mail:        mail_config,
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", "admin"),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 24 * 3, // 3 days
				iss:    "gophersocial",
			},
		},
		redisCfg: redisConfig{
			addr:     env.GetString("REDIS_ADDR", "localhost:6379"),
			password: env.GetString("REDIS_PASSWORD", ""),
			db:       env.GetInt("REDIS_DB", 0),
			enabled:  env.GetBool("REDIS_ENABLED", false),
		},
		rateLimiter: rateLimiterCfg,
	}

	// Cache
	var rdb *redis.Client
	if config.redisCfg.enabled {
		rdb = cache.NewRedisClient(config.redisCfg.addr, config.redisCfg.password, config.redisCfg.db)
		logger.Info("redis connection established")
	}

	jwtAuthenticator := auth.NewJWTAuthenticator(
		config.auth.token.secret,
		config.auth.token.iss,
		config.auth.token.iss,
	)

	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		rateLimiterCfg.RequestPerTimeFrame,
		rateLimiterCfg.TimeFrame,
	)

	app := application{
		config:        config,
		store:         store.NewStore(db),
		cacheStorage:  cache.NewRedisStorage(rdb),
		logger:        logger,
		mail:          mailtrap,
		authenticator: jwtAuthenticator,
		rateLimiter:   rateLimiter,
	}

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
