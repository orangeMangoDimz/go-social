package middlewareHandler

import (
	"context"
	"net/http"

	usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"
)

type MiddlewareProvider interface {
	AuthTokenMiddleware(next http.Handler) http.Handler
	BasicAuthMiddleware() func(http.Handler) http.Handler
	CheckPostOwnership(role string, next http.HandlerFunc) http.HandlerFunc
	GetUser(ctx context.Context, userID int64) (*usersEntity.User, error)
	RateLimiterMiddleware(next http.Handler) http.Handler
}
