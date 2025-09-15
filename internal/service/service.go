package service

import (
	"context"
	"time"

	commentsEntity "github.com/orangeMangoDimz/go-social/internal/entities/comments"
	postsEntity "github.com/orangeMangoDimz/go-social/internal/entities/posts"
	usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"
	"github.com/orangeMangoDimz/go-social/internal/storage/postgres/pagination"
)

type UsersService interface {
	GetById(context.Context, int64) (*usersEntity.User, error)
	GetByEmail(context.Context, string) (*usersEntity.User, error)
	FollowUser(context.Context, int64, int64) error
	CreateAndInvite(context.Context, *usersEntity.User, string, time.Duration) error
	Activate(context.Context, string) error
	Delete(context.Context, int64) error
}

type FollowerService interface {
	Follow(context.Context, int64, int64) error
	Unfollow(context.Context, int64, int64) error
}

type PostsService interface {
	GetUserFeed(context.Context, int64, pagination.PaginatedQuery) ([]postsEntity.Feed, error)
	Create(context.Context, *postsEntity.Post) error
	GetById(context.Context, int64) (*postsEntity.Post, error)
	Update(context.Context, *postsEntity.Post) error
	Delete(context.Context, int64) error
}

type RoleService interface {
	GetByName(context.Context, string) (*usersEntity.Role, error)
}

type CommentService interface {
	Create(context.Context, *commentsEntity.Comment) error
	GetByPostID(context.Context, int64) ([]commentsEntity.Comment, error)
}

type Service struct {
	UsersService    UsersService
	FollowerService FollowerService
	PostService     PostsService
	RoleService     RoleService
	CommentService  CommentService
}
