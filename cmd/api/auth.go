package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"

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
