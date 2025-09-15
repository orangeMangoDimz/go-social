package storage

import (
	"context"
	"database/sql"
	"errors"
	"time"

	commentsEntity "github.com/orangeMangoDimz/go-social/internal/entities/comments"
	postsEntity "github.com/orangeMangoDimz/go-social/internal/entities/posts"
	usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"
	"github.com/orangeMangoDimz/go-social/internal/storage/postgres/pagination"
)

var (
	ErrNotFound          = errors.New("RESOURCE NOT FOUND")
	QueryTimeoutDuration = time.Second * 5
	ErrUniqueViolation   = errors.New("DUPLICATE UNIQUE RECORDS")
	ErrDuplicateEmail    = errors.New("a user with that email already exists")
	ErrDuplicateUsername = errors.New("a user with that username already exists")
)

type Storage struct {
	Posts     PostsRepository
	Users     UsersRepository
	Comments  CommentsRepository
	Followers FollowersRepository
	Roles     RolesRepository
}

type UsersRepository interface {
	GetById(context.Context, int64) (*usersEntity.User, error)
	GetByEmail(context.Context, string) (*usersEntity.User, error)
	Create(context.Context, *sql.Tx, *usersEntity.User) error
	CreateAndInvite(context.Context, *usersEntity.User, string, time.Duration) error
	Activate(context.Context, string) error
	Delete(context.Context, int64) error
}

type PostsRepository interface {
	GetById(context.Context, int64) (*postsEntity.Post, error)
	Delete(context.Context, int64) error
	Create(context.Context, *postsEntity.Post) error
	Update(context.Context, *postsEntity.Post) error
	GetUserFeed(context.Context, int64, pagination.PaginatedQuery) ([]postsEntity.Feed, error)
}

type CommentsRepository interface {
	Create(context.Context, *commentsEntity.Comment) error
	GetByPostID(context.Context, int64) ([]commentsEntity.Comment, error)
}

type FollowersRepository interface {
	Follow(ctx context.Context, followedID, userID int64) error
	Unfollow(ctx context.Context, followedID, userID int64) error
}

type RolesRepository interface {
	GetByName(context.Context, string) (*usersEntity.Role, error)
}

func WithTx(db *sql.DB, ctx context.Context, fn func(*sql.Tx) error) error {
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
