package usersHandler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/orangeMangoDimz/go-social/internal/server/http/protocol"
	"github.com/orangeMangoDimz/go-social/internal/service"
	"github.com/orangeMangoDimz/go-social/internal/storage"
	"go.uber.org/zap"
)

type httpHandler struct {
	userService     service.UsersService
	followerService service.FollowerService
	logger          zap.SugaredLogger
}

func newHTTPHandler(userService service.UsersService, followerService service.FollowerService, logger zap.SugaredLogger) *httpHandler {
	return &httpHandler{
		userService:     userService,
		followerService: followerService,
		logger:          logger,
	}
}

// GetUserHandler godoc
//
//	@Summary		Get user by ID
//	@Description	Get detailed information about a specific user. Requires JWT authentication.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			userID	path		int																	true	"User ID"	example(1)
//	@Success		200		{object}	github_com_orangeMangoDimz_go-social_internal_entities_users.User	"User information"
//	@Failure		401		{object}	map[string]string													"Unauthorized - invalid or missing token"
//	@Failure		404		{object}	map[string]string													"User not found"
//	@Failure		500		{object}	map[string]string													"Internal server error"
//	@Router			/users/{userID} [get]
func (h *httpHandler) GetUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		h.logger.Warn("Error parsing userID")
		protocol.BadRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()
	user, err := h.userService.GetById(ctx, userID)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			h.logger.Warnf("User with user_id %d not found", userID)
			protocol.NotFoundResponse(w, r, err)
		default:
			h.logger.Errorw("Failed to get user data", "user_id", userID, "error", err)
			protocol.InternalServerError(w, r, err)
		}
		return
	}

	if err := protocol.JsonResponse(w, http.StatusOK, user); err != nil {
		h.logger.Errorw("Failed to send response", "error", err)
		protocol.InternalServerError(w, r, err)
		return
	}
}

// FollowUserHandler godoc
//
//	@Summary		Follow a user
//	@Description	Follow another user to see their posts in your feed. Requires JWT authentication.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			userID	path	int	true	"ID of the user to follow"	example(1)
//	@Success		204		"Successfully followed user"
//	@Failure		400		{object}	map[string]string	"Bad request"
//	@Failure		401		{object}	map[string]string	"Unauthorized - invalid or missing token"
//	@Failure		404		{object}	map[string]string	"User not found"
//	@Failure		409		{object}	map[string]string	"Already following this user"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/users/{userID}/follow [put]
func (h *httpHandler) FollowUserHandler(w http.ResponseWriter, r *http.Request) {
	// user to follow
	followerUser := protocol.GetUserFromContext(r)
	if followerUser == nil {
		h.logger.Warnf("User with user_id %d not found", followerUser.ID)
		protocol.NotFoundResponse(w, r, storage.ErrNotFound)
		return
	}

	followedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		h.logger.Warn("Error parsing userID")
		protocol.BadRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	if err := h.followerService.Follow(ctx, followerUser.ID, followedID); err != nil {
		switch {
		case errors.Is(err, storage.ErrUniqueViolation):
			h.logger.Warnw("Conflict on user followers", "error", err)
			protocol.ConflictResponse(w, r, errors.New("YOU HAVE FOLLOWED THIS USER"))
		default:
			h.logger.Errorw("Failed to follow user", "error", err)
			protocol.InternalServerError(w, r, err)
		}
		return
	}

	if err := protocol.JsonResponse(w, http.StatusNoContent, nil); err != nil {
		h.logger.Errorw("Failed to send response", "error", err)
		protocol.InternalServerError(w, r, err)
		return
	}
}

// unfollowUserHandler godoc
//
//	@Summary		Unfollow a user
//	@Description	Unfollow a user to stop seeing their posts in your feed. Requires JWT authentication.
//	@Tags			users
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Param			userID	path	int	true	"ID of the user to unfollow"	example(1)
//	@Success		204		"Successfully unfollowed user"
//	@Failure		400		{object}	map[string]string	"Bad request"
//	@Failure		401		{object}	map[string]string	"Unauthorized - invalid or missing token"
//	@Failure		404		{object}	map[string]string	"User not found"
//	@Failure		409		{object}	map[string]string	"Already unfollowed this user"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/users/{userID}/unfollow [put]
func (h *httpHandler) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	// user to follow
	followerUser := protocol.GetUserFromContext(r)
	if followerUser == nil {
		protocol.NotFoundResponse(w, r, storage.ErrNotFound)
		return
	}

	unfollowedID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)
	if err != nil {
		protocol.BadRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	err = h.followerService.Unfollow(ctx, followerUser.ID, unfollowedID)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrUniqueViolation):
			protocol.ConflictResponse(w, r, errors.New("YOU HAVE UNFOLLOWED THIS USER"))
		default:
			protocol.InternalServerError(w, r, err)
		}
		return
	}

	if err := protocol.JsonResponse(w, http.StatusNoContent, nil); err != nil {
		protocol.InternalServerError(w, r, err)
		return
	}
}

// activateUserHandler godoc
//
//	@Summary		Activate user account
//	@Description	Activate a user account using the activation token received during registration
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			token	path	string	true	"Activation token"	example("550e8400-e29b-41d4-a716-446655440000")
//	@Success		202		"User account activated successfully"
//	@Failure		404		{object}	map[string]string	"Token not found or invalid"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/users/activate/{token} [put]
func (h *httpHandler) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	ctx := r.Context()

	err := h.userService.Activate(ctx, token)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			protocol.NotFoundResponse(w, r, err)
		default:
			protocol.InternalServerError(w, r, err)
		}
		return
	}

	if err := protocol.JsonResponse(w, http.StatusAccepted, nil); err != nil {
		protocol.InternalServerError(w, r, err)
		return
	}
}
