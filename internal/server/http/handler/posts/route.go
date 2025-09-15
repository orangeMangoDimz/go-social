package postsHandler

import (
	"github.com/go-chi/chi/v5"
	middlewareHandler "github.com/orangeMangoDimz/go-social/internal/server/http/middleware"
	"github.com/orangeMangoDimz/go-social/internal/service"
	"go.uber.org/zap"
)

func RegisterRoute(
	middlewareProvider middlewareHandler.MiddlewareProvider,
	postService service.PostsService,
	commentService service.CommentService,
	logger zap.SugaredLogger,
) func(chi.Router) {
	return func(r chi.Router) {
		handler := newHTTPHandler(postService, commentService, logger)
		r.Use(middlewareProvider.AuthTokenMiddleware)
		r.Post("/", handler.createPostHandler)
		r.Route("/{postID}", func(r chi.Router) {
			r.Use(handler.postContextMiddleware)
			r.Get("/", handler.getPostHandler)
			r.Patch("/", middlewareProvider.CheckPostOwnership("moderator", handler.updatePostHandler))
			r.Delete("/", middlewareProvider.CheckPostOwnership("admin", handler.deletePostHandler))
		})
		r.Group(func(r chi.Router) {
			r.Use(middlewareProvider.AuthTokenMiddleware)
			r.Get("/feed", handler.getUserPostFeed)
		})
	}
}
