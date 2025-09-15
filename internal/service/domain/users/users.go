package usersService

import (
	"context"
	"time"

	"github.com/orangeMangoDimz/go-social/internal/config"
	usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"
	"github.com/orangeMangoDimz/go-social/internal/storage"
	"go.uber.org/zap"
)

type UserService struct {
	userRepository storage.UsersRepository
	Logger         *zap.SugaredLogger
	Config         config.Config
}

func NewUserService(userRepository storage.UsersRepository, logger *zap.SugaredLogger, config config.Config) *UserService {
	return &UserService{
		userRepository: userRepository,
		Logger:         logger,
		Config:         config,
	}
}

func (s *UserService) GetById(ctx context.Context, userID int64) (*usersEntity.User, error) {
	user, err := s.userRepository.GetById(ctx, userID)
	return user, err
}

func (s *UserService) GetByEmail(ctx context.Context, userEmail string) (*usersEntity.User, error) {
	user, err := s.userRepository.GetByEmail(ctx, userEmail)
	return user, err
}

func (s *UserService) FollowUser(ctx context.Context, followedUserID int64, followedID int64) error {
	return nil
}

func (s *UserService) CreateAndInvite(ctx context.Context, user *usersEntity.User, token string, invitationExp time.Duration) error {
	err := s.userRepository.CreateAndInvite(ctx, user, token, invitationExp)
	return err
}

func (s *UserService) Activate(ctx context.Context, token string) error {
	err := s.userRepository.Activate(ctx, token)
	return err
}

func (s *UserService) Delete(ctx context.Context, userID int64) error {
	err := s.userRepository.Delete(ctx, userID)
	return err
}
