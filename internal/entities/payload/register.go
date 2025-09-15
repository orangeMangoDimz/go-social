package payloadEntity

import usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"

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
	*usersEntity.User
	Token string `json:"token" example:"550e8400-e29b-41d4-a716-446655440000"` // Activation token
}

type CreateUserTokenPayload struct {
	Email    string `json:"email" validate:"required,email,max=255" example:"johndoe@example.com"` // Email address
	Password string `json:"password" validate:"required,min=3,max=72" example:"securepassword123"` // Password
}
