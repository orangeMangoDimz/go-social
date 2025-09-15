package postsEntity

import (
	commentsEntity "github.com/orangeMangoDimz/go-social/internal/entities/comments"
	usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"
)

// Post represents a social media post
//
//	@Description	Social media post with content, tags and metadata
type Post struct {
	ID        int64                    `json:"id" example:"1"`                            // Post ID
	Content   string                   `json:"content" example:"This is my post content"` // Post content
	Title     string                   `json:"title" example:"My First Post"`             // Post title
	UserId    int64                    `json:"user_id" example:"123"`                     // ID of the user who created the post
	Tags      []string                 `json:"tags" example:"golang,programming"`         // Post tags
	CreatedAt string                   `json:"created_at" example:"2024-01-01 12:00:00"`  // Post creation timestamp
	UpdatedAt string                   `json:"updated_at" example:"2024-01-01 12:30:00"`  // Post last update timestamp
	Version   int                      `json:"version" example:"1"`                       // Post version for optimistic locking
	Comments  []commentsEntity.Comment `json:"comments"`                                  // Comments on this post
	User      usersEntity.User         `json:"user"`                                      // User who created the post
}
