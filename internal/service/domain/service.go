package domain

import (
	"github.com/orangeMangoDimz/go-social/internal/config"
	"github.com/orangeMangoDimz/go-social/internal/service"
	commentsService "github.com/orangeMangoDimz/go-social/internal/service/comments"
	followersService "github.com/orangeMangoDimz/go-social/internal/service/domain/followers"
	postsService "github.com/orangeMangoDimz/go-social/internal/service/domain/posts"
	rolesService "github.com/orangeMangoDimz/go-social/internal/service/domain/roles"
	usersService "github.com/orangeMangoDimz/go-social/internal/service/domain/users"
	"github.com/orangeMangoDimz/go-social/internal/storage"
	"go.uber.org/zap"
)

func NewService(repository storage.Storage, logger *zap.SugaredLogger, config config.Config) *service.Service {
	return &service.Service{
		UsersService:    usersService.NewUserService(repository.Users, logger, config),
		FollowerService: followersService.NewFollowerService(repository.Followers),
		PostService:     postsService.NewPostService(repository.Posts),
		RoleService:     rolesService.NewRoleService(repository.Roles),
		CommentService:  commentsService.NewPostService(repository.Comments),
	}
}
