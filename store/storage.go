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
		Create(context.Context, *User) error
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
