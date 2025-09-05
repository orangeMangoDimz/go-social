package store

import (
	"context"
	"database/sql"
	"errors"
)

// User represents a user in the system
//
//	@Description	User account information
type User struct {
	ID        int64  `json:"id" example:"1"`                           // User ID
	Username  string `json:"username" example:"johndoe"`               // Username
	Email     string `json:"email" example:"johndoe@example.com"`      // Email address
	Password  string `json:"-"`                                        // Password (never returned in responses)
	CreatedAt string `json:"created_at" example:"2024-01-01 12:00:00"` // Account creation timestamp
}

type UserStore struct {
	db *sql.DB
}

func (s *UserStore) GetById(ctx context.Context, userID int64) (*User, error) {
	query := `
		SELECT id, username, email, created_at
		FROM users
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	var user User
	err := s.db.QueryRowContext(
		ctx,
		query,
		userID,
	).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (s *UserStore) Create(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (username, password, email) 
		VALUES ($1, $2, $3) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	err := s.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Password,
		user.Email,
	).Scan(
		&user.ID,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}
