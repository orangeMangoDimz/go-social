package postsService

import (
	"context"

	postsEntity "github.com/orangeMangoDimz/go-social/internal/entities/posts"
	"github.com/orangeMangoDimz/go-social/internal/storage"
	"github.com/orangeMangoDimz/go-social/internal/storage/postgres/pagination"
)

type PostService struct {
	postRepository storage.PostsRepository
	// validation
}

func NewPostService(postRepository storage.PostsRepository) *PostService {
	return &PostService{
		postRepository: postRepository,
		// validation
	}
}

func (s *PostService) GetUserFeed(ctx context.Context, userID int64, fq pagination.PaginatedQuery) ([]postsEntity.Feed, error) {
	feeds, err := s.postRepository.GetUserFeed(ctx, userID, fq)
	return feeds, err
}

func (s *PostService) Create(ctx context.Context, post *postsEntity.Post) error {
	err := s.postRepository.Create(ctx, post)
	return err
}

func (s *PostService) GetById(ctx context.Context, postId int64) (*postsEntity.Post, error) {
	post, err := s.postRepository.GetById(ctx, postId)
	return post, err
}

func (s *PostService) Update(ctx context.Context, post *postsEntity.Post) error {
	err := s.postRepository.Update(ctx, post)
	return err
}

func (s *PostService) Delete(ctx context.Context, postID int64) error {
	err := s.postRepository.Delete(ctx, postID)
	return err
}
