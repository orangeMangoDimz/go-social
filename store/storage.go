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
)

type Storage struct {
	Posts interface {
		GetById(context.Context, int64) (*Post, error)
		Delete(context.Context, int64) error
		Create(context.Context, *Post) error
		Update(context.Context, *Post) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
	Comments interface {
		Create(context.Context, *Comment) error
		GetByPostID(context.Context, int64) ([]Comment, error)
	}
}

func NewStore(db *sql.DB) Storage {
	return Storage{
		Posts:    &PostStore{db: db},
		Users:    &UserStore{db: db},
		Comments: &CommentStore{db: db},
	}
}
