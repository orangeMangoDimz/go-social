package httpserver

import (
	"github.com/orangeMangoDimz/go-social/internal/config"
	"github.com/orangeMangoDimz/go-social/internal/mailer"
	"github.com/orangeMangoDimz/go-social/internal/ratelimiter"
	authHandler "github.com/orangeMangoDimz/go-social/internal/server/http/handler/auth"
	"github.com/orangeMangoDimz/go-social/internal/service"
	"github.com/orangeMangoDimz/go-social/internal/storage"
	"github.com/orangeMangoDimz/go-social/internal/storage/cache"
	"go.uber.org/zap"
)

type Application struct {
	Config        config.Config
	Store         storage.Storage
	CacheStorage  cache.Storage
	Logger        *zap.SugaredLogger
	Mail          mailer.Client
	Authenticator authHandler.Authenticator
	RateLimiter   ratelimiter.Limiter
	Us            service.UsersService
	Services      service.Service
}
