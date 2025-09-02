package store

import (
	"context"
	"database/sql"
)

type Storage struct {
	Posts interface {
		Create(context.Context) error
	}
	Users interface {
		Create(context.Context, *User) error
	}
}

func NewStore(db *sql.DB) Storage {
	return Storage{
		Posts: &PostStore{db: db},
		Users: &UserStore{db: db},
	}
}
