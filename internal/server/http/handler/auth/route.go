package authHandler

import (
	"github.com/go-chi/chi/v5"
	"github.com/orangeMangoDimz/go-social/internal/config"
	"github.com/orangeMangoDimz/go-social/internal/mailer"
	"github.com/orangeMangoDimz/go-social/internal/service"
	"go.uber.org/zap"
)

func RegisterRoute(userService service.UsersService, logger zap.SugaredLogger, mailer mailer.Client, config config.Config, authenticator Authenticator) func(chi.Router) {
	return func(r chi.Router) {
		handler := newHTTPHandler(userService, logger, mailer, config, authenticator)
		r.Post("/user", handler.registerUserHandler)
		// Add other auth routes here as needed
		r.Post("/token", handler.createTokenHandler)
	}
}
