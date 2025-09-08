package main

import (
	"context"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/orangeMangoDimz/go-social/store"
)

type userKey string

const userCtx userKey = "user"

// getUserHandler retrieves a user by ID
//
//	@Summary		Get user by ID
//	@Description	Get detailed information about a specific user. Requires JWT authentication.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int					true	"User ID"	example(1)
//	@Success		200		{object}	store.User			"User information"
//	@Failure		401		{object}	map[string]string	"Unauthorized - invalid or missing token"
//	@Failure		404		{object}	map[string]string	"User not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/users/{userID} [get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)
	if user == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	if err := app.jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// followUserHandler allows a user to follow another user
//
//	@Summary		Follow a user
//	@Description	Follow another user to see their posts in your feed. Requires JWT authentication.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"ID of the user to follow"	example(1)
//	@Success		204		"Successfully followed user"
//	@Failure		400		{object}	map[string]string	"Bad request"
//	@Failure		401		{object}	map[string]string	"Unauthorized - invalid or missing token"
//	@Failure		404		{object}	map[string]string	"User not found"
//	@Failure		409		{object}	map[string]string	"Already following this user"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	// user to follow
	followerUser := getUserFromContext(r)
	if followerUser == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, followerUser.ID, followedID); err != nil {
		switch {
		case errors.Is(err, store.ErrUniqueViolation):
			app.conflictResponse(w, r, errors.New("YOU HAVE FOLLOWED THIS USER"))
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// unfollowUserHandler allows a user to unfollow another user
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user to stop seeing their posts in your feed. Requires JWT authentication.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int	true	"ID of the user to unfollow"	example(1)
//	@Success		204		"Successfully unfollowed user"
//	@Failure		400		{object}	map[string]string	"Bad request"
//	@Failure		401		{object}	map[string]string	"Unauthorized - invalid or missing token"
//	@Failure		404		{object}	map[string]string	"User not found"
//	@Failure		409		{object}	map[string]string	"Already unfollowed this user"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Security		BearerAuth
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	// user to follow
	followerUser := getUserFromContext(r)
	if followerUser == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	unfollowedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Unfollow(ctx, followerUser.ID, unfollowedID); err != nil {
		switch {
		case errors.Is(err, store.ErrUniqueViolation):
			app.conflictResponse(w, r, errors.New("YOU HAVE UNFOLLOWED THIS USER"))
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}

// activateUserHandler activates a user account using the provided token
//
//	@Summary		Activate user account
//	@Description	Activate a user account using the activation token received during registration
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			token	path	string	true	"Activation token"	example("550e8400-e29b-41d4-a716-446655440000")
//	@Success		204		"User account activated successfully"
//	@Failure		404		{object}	map[string]string	"Token not found or invalid"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	ctx := r.Context()

	if err := app.store.Users.Activate(ctx, token); err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := app.jsonResponse(w, http.StatusAccepted, nil); err != nil {
		app.internalServerError(w, r, err)
		return
	}

}

func (app *application) userContextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idParm := chi.URLParam(r, "userID")
		id, err := strconv.ParseInt(idParm, 10, 64)
		if err != nil {
			app.internalServerError(w, r, err)
			return
		}
		ctx := r.Context()

		user, err := app.store.Users.GetById(ctx, id)
		if err != nil {
			switch {
			case errors.Is(err, store.ErrNotFound):
				app.notFoundResponse(w, r, err)
			default:
				app.internalServerError(w, r, err)
			}
			return
		}

		ctx = context.WithValue(ctx, userCtx, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func getUserFromContext(r *http.Request) *store.User {
	user, ok := r.Context().Value(userCtx).(*store.User)
	if !ok {
		return nil
	}
	return user
}
