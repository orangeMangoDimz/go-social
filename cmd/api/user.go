package main

import (
	"context"
	"errors"
	"fmt"
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
//	@Description	Get detailed information about a specific user
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path		int					true	"User ID"	example(1)
//	@Success		200		{object}	store.User			"User information"
//	@Failure		404		{object}	map[string]string	"User not found"
//	@Failure		500		{object}	map[string]string	"Internal server error"
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

// FollowUser represents the request payload for following/unfollowing a user
//
//	@Description	Request payload for follow/unfollow operations
type FollowUser struct {
	UserID int64 `json:"user_id" example:"123" validate:"required"` // User ID to follow/unfollow
}

// followUserHandler allows a user to follow another user
//
//	@Summary		Follow a user
//	@Description	Follow another user to see their posts in your feed
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int			true	"ID of the user doing the following"	example(1)
//	@Param			payload	body	FollowUser	true	"User to follow"
//	@Success		204		"Successfully followed user"
//	@Failure		400		{object}	map[string]string	"Bad request"
//	@Failure		404		{object}	map[string]string	"User not found"
//	@Failure		409		{object}	map[string]string	"Already following this user"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/users/{userID}/follow [put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("test")
	// user to follow
	followerUser := getUserFromContext(r)
	if followerUser == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	// Revert from ctx
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Follow(ctx, followerUser.ID, payload.UserID); err != nil {
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
//	@Description	Unfollow a user to stop seeing their posts in your feed
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Param			userID	path	int			true	"ID of the user doing the unfollowing"	example(1)
//	@Param			payload	body	FollowUser	true	"User to unfollow"
//	@Success		204		"Successfully unfollowed user"
//	@Failure		400		{object}	map[string]string	"Bad request"
//	@Failure		404		{object}	map[string]string	"User not found"
//	@Failure		409		{object}	map[string]string	"Already unfollowed this user"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/users/{userID}/unfollow [put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	// user to follow
	unfollowedUser := getUserFromContext(r)
	if unfollowedUser == nil {
		app.notFoundResponse(w, r, store.ErrNotFound)
		return
	}

	// Revert from ctx
	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := app.store.Followers.Unfollow(ctx, unfollowedUser.ID, payload.UserID); err != nil {
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
