package usersHandler

import (
	"github.com/go-chi/chi/v5"
	middlewareHandler "github.com/orangeMangoDimz/go-social/internal/server/http/middleware"
	"github.com/orangeMangoDimz/go-social/internal/service"
	"go.uber.org/zap"
)

func RegisterRoute(
	middlewareProvider middlewareHandler.MiddlewareProvider,
	userService service.UsersService,
	followerService service.FollowerService,
	logger zap.SugaredLogger,
) func(chi.Router) {
	return func(r chi.Router) {
		handler := newHTTPHandler(userService, followerService, logger)
		r.Put("/activate/{token}", handler.activateUserHandler)
		r.Route("/{userID}", func(r chi.Router) {
			r.Use(middlewareProvider.AuthTokenMiddleware)
			r.Get("/", handler.GetUserHandler)
			r.Put("/follow", handler.FollowUserHandler)
			r.Put("/unfollow", handler.unfollowUserHandler)
		})

	}
}
