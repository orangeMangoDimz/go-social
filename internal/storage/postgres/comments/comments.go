package comments

import (
	"context"
	"database/sql"

	commentsEntity "github.com/orangeMangoDimz/go-social/internal/entities/comments"
	usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"
	"github.com/orangeMangoDimz/go-social/internal/storage"
)

type CommentStore struct {
	Db *sql.DB
}

func (s *CommentStore) Create(ctx context.Context, comment *commentsEntity.Comment) error {
	query := `
		INSERT INTO comments (post_id, user_id, content) 
		VALUES ($1, $2, $3) RETURNING id, created_at
	`

	ctx, cancel := context.WithTimeout(ctx, storage.QueryTimeoutDuration)
	defer cancel()

	err := s.Db.QueryRowContext(
		ctx,
		query,
		comment.PostID,
		comment.UserID,
		comment.Content,
	).Scan(
		&comment.ID,
		&comment.CreatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *CommentStore) GetByPostID(ctx context.Context, postID int64) ([]commentsEntity.Comment, error) {
	query := `
		SELECT c.id, c.post_id, c.user_id, c.content, c.created_at, u.username, u.id
		FROM comments c
		JOIN users u ON u.id = c.user_id
		WHERE c.post_id = $1
		ORDER BY c.created_at DESC
	`

	ctx, cancel := context.WithTimeout(ctx, storage.QueryTimeoutDuration)
	defer cancel()

	rows, err := s.Db.QueryContext(
		ctx,
		query,
		postID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	comments := []commentsEntity.Comment{}

	for rows.Next() {
		var c commentsEntity.Comment
		c.User = usersEntity.User{}
		err := rows.Scan(&c.ID, &c.PostID, &c.UserID, &c.Content, &c.CreatedAt, &c.User.Username, &c.User.ID)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
