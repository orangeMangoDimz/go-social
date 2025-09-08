package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/orangeMangoDimz/go-social/internal/mailer"
	"github.com/orangeMangoDimz/go-social/store"
)

// RegisterUserPayload represents the request payload for user registration
//
//	@Description	Request payload for user registration
type RegisterUserPayload struct {
	Username string `json:"username" validate:"required,max=100" example:"johndoe"`                // Username (max 100 characters)
	Email    string `json:"email" validate:"required,email,max=255" example:"johndoe@example.com"` // Email address (max 255 characters)
	Password string `json:"password" validate:"required,min=3,max=72" example:"securepassword123"` // Password (3-72 characters)
}

// UserWithToken represents a user with an activation token
//
//	@Description	User information with activation token
type UserWithToken struct {
	*store.User
	Token string `json:"token" example:"550e8400-e29b-41d4-a716-446655440000"` // Activation token
}

// registerUserHandler creates a new user account and sends an activation token
//
//	@Summary		Register a new user
//	@Description	Create a new user account with username, email and password. Returns user information with an activation token.
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		RegisterUserPayload	true	"User registration data"
//	@Success		202		{object}	UserWithToken		"User created successfully, activation required"
//	@Failure		400		{object}	map[string]string	"Bad request - validation error or duplicate email/username"
//	@Failure		500		{object}	map[string]string	"Internal server error"
//	@Router			/authentication/user [post]
func (app *application) registerUserHandler(w http.ResponseWriter, r *http.Request) {
	var payload RegisterUserPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	user := &store.User{
		Username: payload.Username,
		Email:    payload.Email,
	}

	// Hash user password
	if err := user.Password.Set(payload.Password); err != nil {
		app.internalServerError(w, r, err)
		return
	}
	ctx := r.Context()

	plainToken := uuid.New().String()

	// Hash the token
	hash := sha256.Sum256([]byte(plainToken))
	hashToken := hex.EncodeToString(hash[:])

	// Store the user
	if err := app.store.Users.CreateAndInvite(ctx, user, hashToken, app.config.mail.exp); err != nil {
		switch {
		case errors.Is(err, store.ErrDuplicateEmail):
			app.badRequestResponse(w, r, err)
		case errors.Is(err, store.ErrDuplicateUsername):
			app.badRequestResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// send the email with plain token
	userWithToken := &UserWithToken{
		User:  user,
		Token: plainToken,
	}

	activationURL := fmt.Sprintf("%s/confirm/%s", app.config.frontendURL, plainToken)

	isProdEnv := app.config.env == "production"
	vars := struct {
		Username      string
		ActivationURL string
	}{
		Username:      user.Username,
		ActivationURL: activationURL,
	}

	// Send email
	status, err := app.mail.Send(mailer.UserWelcomeTemplate, user.Username, user.Email, vars, !isProdEnv)
	if err != nil {
		app.logger.Errorw("error sending welcome email", "error", err)

		// rollback user creation if email fails (SAGA pattern)
		if err := app.store.Users.Delete(ctx, user.ID); err != nil {
			app.logger.Errorw("error deleting user", "error", err)
		}

		app.internalServerError(w, r, err)
		return
	}

	app.logger.Infow("Email sent", "status code", status)

	if err := app.jsonResponse(w, http.StatusCreated, userWithToken); err != nil {
		app.internalServerError(w, r, err)
	}

}

// CreateUserTokenPayload represents the request payload for creating an authentication token
//
//	@Description	Request payload for creating a JWT authentication token
type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255" example:"johndoe@example.com"` // Email address
	Password string `json:"password" validate:"required,min=3,max=72" example:"securepassword123"` // Password
}

// TokenResponse represents the response containing a JWT token
//
//	@Description	Response containing a JWT authentication token
type TokenResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."` // JWT authentication token
}

// createTokenHandler creates a JWT authentication token for user login
//
//	@Summary		Login and get authentication token
//	@Description	Authenticate user with email and password, returns a JWT token for API access
//	@Tags			authentication
//	@Accept			json
//	@Produce		json
//	@Param			payload	body		CreateUserTokenPayload	true	"User login credentials"
//	@Success		201		{object}	TokenResponse			"JWT token created successfully"
//	@Failure		400		{object}	map[string]string		"Bad request - validation error"
//	@Failure		401		{object}	map[string]string		"Unauthorized - invalid credentials"
//	@Failure		500		{object}	map[string]string		"Internal server error"
//	@Router			/authentication/token [post]
func (app *application) createTokenHandler(w http.ResponseWriter, r *http.Request) {
	// parse payload credentials
	var payload CreateUserTokenPayload
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := validate.Struct(payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	// fetch the user (check if the user exists) form the payload
	ctx := r.Context()
	user, err := app.store.Users.GetByEmail(ctx, payload.Email)
	if err != nil {
		switch {
		case errors.Is(err, store.ErrNotFound):
			app.unauthorizedErrorResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	// verify the password
	if err := user.Password.Compare(payload.Password); err != nil {
		app.unauthorizedErrorResponse(w, r, err)
		return
	}

	// generate the token -> add claims
	claims := jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(app.config.auth.token.exp).Unix(),
		"iat": time.Now().Unix(),
		"nbf": time.Now().Unix(),
		"iss": app.config.auth.token.iss,
		"aud": app.config.auth.token.iss,
	}
	token, err := app.authenticator.GenerateToken(claims)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	// send it to the client
	response := TokenResponse{Token: token}
	if err := app.jsonResponse(w, http.StatusCreated, response); err != nil {
		app.internalServerError(w, r, err)
	}
}
