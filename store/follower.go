package store

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
)

type Follower struct {
	UserID     int64  `json:"user_id"`
	FollowerID int64  `json:"follower_id"`
	CreatedAT  string `json:"created_at"`
}

type FollowerStore struct {
	db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followedID, userID int64) error {
	query := `
	INSERT INTO followers (user_id, follower_id) 
	VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, QueryTimeoutDuration)
	defer cancel()

	_, err := s.db.ExecContext(ctx, query, userID, followedID)
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.As(err, &pqErr):
			if pqErr.Code == "23505" {
				return ErrUniqueViolation
			}
		}
		return err
	}

	return err
}

