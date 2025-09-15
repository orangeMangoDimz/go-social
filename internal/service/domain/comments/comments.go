package commentsService

import (
	"context"

	commentsEntity "github.com/orangeMangoDimz/go-social/internal/entities/comments"
	"github.com/orangeMangoDimz/go-social/internal/storage"
)

type CommentService struct {
	commentRepository storage.CommentsRepository
	// validation
}

func NewPostService(commentRepository storage.CommentsRepository) *CommentService {
	return &CommentService{
		commentRepository: commentRepository,
		// validation
	}
}

func (s *CommentService) Create(ctx context.Context, comment *commentsEntity.Comment) error {
	err := s.commentRepository.Create(ctx, comment)
	return err
}

func (s *CommentService) GetByPostID(ctx context.Context, postID int64) ([]commentsEntity.Comment, error) {
	comment, err := s.commentRepository.GetByPostID(ctx, postID)
	return comment, err
}
