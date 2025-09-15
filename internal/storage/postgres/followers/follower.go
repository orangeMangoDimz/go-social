package followers

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	"github.com/orangeMangoDimz/go-social/internal/storage"
)

type FollowerStore struct {
	Db *sql.DB
}

func (s *FollowerStore) Follow(ctx context.Context, followedID, userID int64) error {
	query := `
	INSERT INTO followers (user_id, follower_id) 
	VALUES ($1, $2)
	`

	ctx, cancel := context.WithTimeout(ctx, storage.QueryTimeoutDuration)
	defer cancel()

	_, err := s.Db.ExecContext(ctx, query, userID, followedID)
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.As(err, &pqErr):
			if pqErr.Code == "23505" {
				return storage.ErrUniqueViolation
			}
		}
		return err
	}

	return err
}

func (s *FollowerStore) Unfollow(ctx context.Context, followedID, userID int64) error {
	query := `
	DELETE FROM followers
	WHERE user_id = $1 AND follower_id = $2
	`

	ctx, cancel := context.WithTimeout(ctx, storage.QueryTimeoutDuration)
	defer cancel()

	_, err := s.Db.ExecContext(ctx, query, userID, followedID)
	if err != nil {
		var pqErr *pq.Error
		switch {
		case errors.As(err, &pqErr):
			if pqErr.Code == "23505" {
				return storage.ErrUniqueViolation
			}
		}
		return err
	}

	return err
}
