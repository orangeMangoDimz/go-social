package commentsEntity

import usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"

// Comment represents a comment on a post
//
//	@Description	Comment on a social media post
type Comment struct {
	ID        int64            `json:"id" example:"1"`                           // Comment ID
	PostID    int64            `json:"post_id" example:"123"`                    // ID of the post this comment belongs to
	UserID    int64            `json:"user_id" example:"456"`                    // ID of the user who made the comment
	Content   string           `json:"content" example:"Great post!"`            // Comment content
	CreatedAt string           `json:"created_at" example:"2024-01-01 12:00:00"` // Comment creation timestamp
	User      usersEntity.User `json:"user"`                                     // User who made the comment
}
