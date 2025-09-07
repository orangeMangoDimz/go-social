package store

import (
	"context"
	"database/sql"
	"errors"
	"time"
)

var (
	ErrNotFound          = errors.New("RESOURCE NOT FOUND")
	QueryTimeoutDuration = time.Second * 5
	ErrUniqueViolation   = errors.New("DUPLICATE UNIQUE RECORDS")
	InitSinceTime        = "2025-01-01 00:00:00"
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type Storage struct {
	Posts interface {
		GetById(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Create(context.Context, *Post) error
		Update(context.Context, *Post) error
		GetUserFeed(context.Context, int64, PaginatedQuery) ([]Feed, error)
	}
	Users interface {
		GetById(context.Context, int64) (*User, error)
		Create(context.Context, *sql.Tx, *User) error
		CreateAndInvite(context.Context, *User, string, time.Duration) error
		Activate(context.Context, string) error
		Delete(context.Context, int64) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(context.Context, int64) ([]Comment, error)
	}
	Followers interface {
		Follow(ctx context.Context, followedID, userID int64) error
		Unfollow(ctx context.Context, followedID, userID int64) error
	}
}

func NewStore(db *sql.DB) Storage {
	return Storage{
		Posts:     &PostStore{db: db},
		Users:     &UserStore{db: db},
		Comments:  &CommentStore{db: db},
		Followers: &FollowerStore{db: db},
	}
}

func withTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}

	return tx.Commit()
}
