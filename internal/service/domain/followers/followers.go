package followersService

import (
	"context"

	"github.com/orangeMangoDimz/go-social/internal/storage"
)

type FollowerService struct {
	followerRepository storage.FollowersRepository
	// validation
}

func NewFollowerService(followerRepository storage.FollowersRepository) *FollowerService {
	return &FollowerService{
		followerRepository: followerRepository,
		// validation
	}
}

func (s *FollowerService) Follow(ctx context.Context, followedID, userID int64) error {
	err := s.followerRepository.Follow(ctx, followedID, userID)
	return err
}

func (s *FollowerService) Unfollow(ctx context.Context, followedID, userID int64) error {
	err := s.followerRepository.Unfollow(ctx, followedID, userID)
	return err
}
