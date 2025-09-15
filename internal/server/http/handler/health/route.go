package healthHandler

import (
	"github.com/go-chi/chi/v5"
	"github.com/orangeMangoDimz/go-social/internal/config"
	middlewareHandler "github.com/orangeMangoDimz/go-social/internal/server/http/middleware"
)

func RegisterRoute(middlewareProvider middlewareHandler.MiddlewareProvider, config config.Config, version string) func(chi.Router) {
	return func(r chi.Router) {
		handler := newHTTPHandler(config, version)
		r.With(middlewareProvider.BasicAuthMiddleware()).Get("/health", handler.healthCheckHandler)
	}
}
