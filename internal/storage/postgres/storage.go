package postgres

import (
	"database/sql"

	"github.com/orangeMangoDimz/go-social/internal/storage"
	"github.com/orangeMangoDimz/go-social/internal/storage/postgres/comments"
	"github.com/orangeMangoDimz/go-social/internal/storage/postgres/followers"
	"github.com/orangeMangoDimz/go-social/internal/storage/postgres/posts"
	"github.com/orangeMangoDimz/go-social/internal/storage/postgres/roles"
	"github.com/orangeMangoDimz/go-social/internal/storage/postgres/users"
)

func NewStore(db *sql.DB) storage.Storage {
	return storage.Storage{
		Posts:     &posts.PostStore{Db: db},
		Users:     &users.UserStore{Db: db},
		Comments:  &comments.CommentStore{Db: db},
		Followers: &followers.FollowerStore{Db: db},
		Roles:     &roles.RoleStore{Db: db},
	}
}
