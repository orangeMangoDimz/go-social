package authHandler

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/orangeMangoDimz/go-social/internal/config"
	payloadEntity "github.com/orangeMangoDimz/go-social/internal/entities/payload"
	usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"
	"github.com/orangeMangoDimz/go-social/internal/mailer"
	"github.com/orangeMangoDimz/go-social/internal/server/http/protocol"
	"github.com/orangeMangoDimz/go-social/internal/service"
	"github.com/orangeMangoDimz/go-social/internal/storage"
	"go.uber.org/zap"
)

type httpHandler struct {
	userService   service.UsersService
	authenticator Authenticator
	logger        zap.SugaredLogger
	mailer        mailer.Client
	config        config.Config
}

func newHTTPHandler(userService service.UsersService, logger zap.SugaredLogger, mailer mailer.Client, config config.Config, authenticator Authenticator) *httpHandler {
	return &httpHandler{
		userService:   userService,
		authenticator: authenticator,
		logger:        logger,
		mailer:        mailer,
		config:        config,
	}
}

// registerUserHandler godoc
//
//	@Summary		Register a new user
//	@Description	Create a new user account with username, email and password. Returns user information with an activation token.
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		github_com_orangeMangoDimz_go-social_internal_entities_payload.RegisterUserPayload	true	"User registration data"
//	@Success		201		{object}	github_com_orangeMangoDimz_go-social_internal_entities_payload.UserWithToken		"User created successfully, activation required"
//	@Failure		400		{object}	map[string]string																	"Bad request - validation error or duplicate email/username"
//	@Failure		500		{object}	map[string]string																	"Internal server error"
//	@Router			/authentication/user [post]
func (h *httpHandler) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload payloadEntity.RegisterUserPayload
	if err := protocol.ReadJSON(w, r, &payload); err != nil {
		protocol.BadRequestResponse(w, r, err)
		return
	}

	if err := protocol.ValidateStruct(payload); err != nil {
		protocol.BadRequestResponse(w, r, err)
		return
	}

	user := &usersEntity.User{
		Username: payload.Username,
		Email:    payload.Email,
		Role: usersEntity.Role{
			Name: "user",
		},
	}

	// Hash user password
	if err := user.Password.Set(payload.Password); err != nil {
		protocol.InternalServerError(w, r, err)
		return
	}
	ctx := r.Context()

	plainToken := uuid.New().String()

	// Hash the token
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	// Store the user
	err := h.userService.CreateAndInvite(ctx, user, hashToken, h.config.Mail.Exp)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrDuplicateEmail):
			protocol.BadRequestResponse(w, r, err)
		case errors.Is(err, storage.ErrDuplicateUsername):
			protocol.BadRequestResponse(w, r, err)
		default:
			protocol.InternalServerError(w, r, err)
		}
		return
	}

	// send the email with plain token
	userWithToken := &payloadEntity.UserWithToken{
		User:  user,
		Token: plainToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", h.config.FrontendURL, plainToken)

	isProdEnv := h.config.Env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	// Send email
	time.Sleep(time.Second * 5)
	status, err := h.mailer.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		h.logger.Errorw("error sending welcome email", "error", err)

		// rollback user creation if email fails (SAGA pattern)
		err := h.userService.Delete(ctx, user.ID)
		if err != nil {
			h.logger.Errorw("error deleting user", "error", err)
		}

		protocol.InternalServerError(w, r, err)
		return
	}

	h.logger.Infow("Email sent", "status code", status)

	if err := protocol.JsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		protocol.InternalServerError(w, r, err)
	}

}

// createTokenHandler godoc
//
//	@Summary		Login and get authentication token
//	@Description	Authenticate user with email and password, returns a JWT token for API access
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		github_com_orangeMangoDimz_go-social_internal_entities_payload.CreateUserTokenPayload	true	"User login credentials"
//	@Success		201		{object}	github_com_orangeMangoDimz_go-social_internal_entities_payload.TokenResponse			"JWT token created successfully"
//	@Failure		400		{object}	map[string]string																		"Bad request - validation error"
//	@Failure		401		{object}	map[string]string																		"Unauthorized - invalid credentials"
//	@Failure		500		{object}	map[string]string																		"Internal server error"
//	@Router			/authentication/token [post]
func (h *httpHandler) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// parse payload credentials
	var payload payloadEntity.CreateUserTokenPayload
	if err := protocol.ReadJSON(w, r, &payload); err != nil {
		protocol.BadRequestResponse(w, r, err)
		return
	}

	var validate = validator.New()
	if err := validate.Struct(payload); err != nil {
		protocol.BadRequestResponse(w, r, err)
		return
	}

	// fetch the user (check if the user exists) form the payload
	ctx := r.Context()
	user, err := h.userService.GetByEmail(ctx, payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, storage.ErrNotFound):
			protocol.UnauthorizedErrorResponse(w, r, err)
		default:
			protocol.InternalServerError(w, r, err)
		}
		return
	}

	// verify the password
	if err := user.Password.Compare(payload.Password); err != nil {
		protocol.UnauthorizedErrorResponse(w, r, err)
		return
	}

	// generate the token -> add claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(h.config.Auth.Token.Exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": h.config.Auth.Token.Iss,
		"aud": h.config.Auth.Token.Iss,
	}
	token, err := h.authenticator.GenerateToken(claims)
	if err != nil {
		protocol.InternalServerError(w, r, err)
		return
	}

	// send it to the client
	response := payloadEntity.TokenResponse{Token: token}
	if err := protocol.JsonResponse(w, http.StatusCreated, response); err != nil {
		protocol.InternalServerError(w, r, err)
	}
}
