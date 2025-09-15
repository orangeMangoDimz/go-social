package httpserver

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt"
	usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"
	"github.com/orangeMangoDimz/go-social/internal/server/http/protocol"
	"github.com/orangeMangoDimz/go-social/internal/storage"
)

func (app *Application) AuthTokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			protocol.UnauthorizedErrorResponse(w, r, fmt.Errorf("bearer authorization is missing"))
			return
		}

		// Expect 'Basic <base64>'
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			protocol.UnauthorizedErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
			return
		}

		token := parts[1]
		jwtToken, err := app.Authenticator.ValidateToken(token)
		if err != nil {
			protocol.UnauthorizedErrorResponse(w, r, err)
			return
		}

		claims := jwtToken.Claims.(jwt.MapClaims)
		userID, err := strconv.ParseInt(fmt.Sprintf("%.f", claims["sub"]), 10, 64)
		if err != nil {
			protocol.UnauthorizedErrorResponse(w, r, err)
			return
		}

		ctx := r.Context()
		user, err := app.Services.UsersService.GetById(ctx, userID)
		if err != nil {
			protocol.InternalServerError(w, r, err)
			return
		}
		ctx = context.WithValue(ctx, protocol.UserCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))

	})
}

// BasicAuthMiddleware applies HTTP Basic auth to protected endpoints.
func (app *Application) BasicAuthMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Get Authorization header
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				protocol.UnauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization is missing"))
				return
			}

			// Expect 'Basic <base64>'
			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Basic" {
				protocol.UnauthorizedBasicErrorResponse(w, r, fmt.Errorf("authorization header is malformed"))
				return
			}

			// Decode base64 credentials
			decoded, err := base64.StdEncoding.DecodeString(parts[1])
			if err != nil {
				protocol.UnauthorizedBasicErrorResponse(w, r, err)
				return
			}

			// Verify username and password
			username := app.Config.Auth.Basic.User
			pass := app.Config.Auth.Basic.Pass

			creds := strings.SplitN(string(decoded), ":", 2)
			if len(creds) != 2 || creds[0] != username || creds[1] != pass {
				protocol.UnauthorizedBasicErrorResponse(w, r, fmt.Errorf("invalid credentials"))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

func (app *Application) CheckPostOwnership(role string, next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := protocol.GetUserFromContext(r)
		post := protocol.GetPostFromContext(r)

		// check if it is user post
		if post.UserId == user.ID {
			next.ServeHTTP(w, r)
			return
		}

		// role precedence check
		ctx := r.Context()
		allowed, err := app.checkRolePrecedence(ctx, user, role)
		if err != nil {
			protocol.InternalServerError(w, r, err)
			return
		}

		if !allowed {
			protocol.ForbiddenResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *Application) checkRolePrecedence(ctx context.Context, user *usersEntity.User, roleName string) (bool, error) {
	role, err := app.Services.RoleService.GetByName(ctx, roleName)
	if err != nil {
		return false, err
	}

	return user.Role.Level >= role.Level, nil
}

func (app *Application) GetUser(ctx context.Context, userID int64) (*usersEntity.User, error) {
	app.Logger.Infow("checking cache for user", "id", userID)
	user, err := app.CacheStorage.Users.Get(ctx, userID)
	if err != nil {
		return nil, err
	}

	if user == nil {
		app.Logger.Infow("cache miss, fetching from DB", "id", userID)
		user, err := app.Services.UsersService.GetById(ctx, userID)
		if err != nil {
			switch {
			case errors.Is(err, storage.ErrNotFound):
				return nil, storage.ErrNotFound
			default:
				return nil, err
			}
		}

		if err := app.CacheStorage.Users.Set(ctx, user); err != nil {
			return nil, err
		}
	} else {
		app.Logger.Infow("cache hit for user", "id", userID)
	}

	return user, nil
}

func (app *Application) RateLimiterMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.Config.RateLimiter.Enabled {
			if allow, retryAfter := app.RateLimiter.Allow(r.RemoteAddr); !allow {
				protocol.RateLimitExceededResponse(w, r, retryAfter.String())
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}
