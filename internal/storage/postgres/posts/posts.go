package posts

import (
	"context"
	"database/sql"
	"errors"

	"github.com/lib/pq"
	postsEntity "github.com/orangeMangoDimz/go-social/internal/entities/posts"
	"github.com/orangeMangoDimz/go-social/internal/storage"
	"github.com/orangeMangoDimz/go-social/internal/storage/postgres/pagination"
)

type PostStore struct {
	Db *sql.DB
}

func (s *PostStore) GetUserFeed(ctx context.Context, userID int64, fq pagination.PaginatedQuery) ([]postsEntity.Feed, error) {
	query := `
		SELECT
			p.id, p.user_id, p.title, p.content, p.created_at, p.tags,
			u.username,
			COUNT(c.id) AS comments_count
		FROM posts p
		LEFT JOIN comments c ON c.post_id = p.id
		LEFT JOIN users u ON p.user_id = u.id
		JOIN followers f ON f.follower_id = p.user_id OR p.user_id = $1
		WHERE 
			(f.user_id = $1 OR p.user_id = $1) AND
			(p.title ILIKE '%' || $4 || '%' OR p.content ILIKE '%' || $4 || '%') AND
			(p.tags @> $5 OR $5 = '{}') AND
			(p.created_at >= $6 AND p.created_at <= $7)
		GROUP BY p.id, u.username
		ORDER BY p.created_at ` + fq.Sort + `
		LIMIT $2
		OFFSET $3
	`

	ctx, cancel := context.WithTimeout(ctx, storage.QueryTimeoutDuration)
	defer cancel()

	rows, err := s.Db.QueryContext(
		ctx,
		query,
		userID,
		fq.Limit,
		fq.Offset,
		fq.Search,
		pq.Array(fq.Tags),
		fq.Since,
		fq.Until,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	feeds := []postsEntity.Feed{}

	for rows.Next() {
		var f postsEntity.Feed
		err := rows.Scan(
			&f.ID,
			&f.UserId,
			&f.Title,
			&f.Content,
			&f.CreatedAt,
			pq.Array(&f.Tags),
			&f.User.Username,
			&f.TotalComments,
		)
		if err != nil {
			return nil, err
		}
		feeds = append(feeds, f)
	}
	return feeds, nil
}

func (s *PostStore) Create(ctx context.Context, post *postsEntity.Post) error {
	query := `
		INSERT INTO posts (content, title, user_id, tags) 
		VALUES ($1, $2, $3, $4) RETURNING id, created_at, updated_at
	`

	ctx, cancel := context.WithTimeout(ctx, storage.QueryTimeoutDuration)
	defer cancel()

	err := s.Db.QueryRowContext(
		ctx,
		query,
		post.Content,
		post.Title,
		post.UserId,
		pq.Array(post.Tags),
	).Scan(
		&post.ID,
		&post.CreatedAt,
		&post.UpdatedAt,
	)

	if err != nil {
		return err
	}
	return nil
}

func (s *PostStore) GetById(ctx context.Context, postId int64) (*postsEntity.Post, error) {
	query := `
		SELECT id, user_id, title, content, tags, created_at, updated_at, version
		FROM posts
		WHERE id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, storage.QueryTimeoutDuration)
	defer cancel()

	var post postsEntity.Post
	err := s.Db.QueryRowContext(
		ctx,
		query,
		postId,
	).Scan(
		&post.ID,
		&post.UserId,
		&post.Title,
		&post.Content,
		pq.Array(&post.Tags),
		&post.CreatedAt,
		&post.UpdatedAt,
		&post.Version,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, storage.ErrNotFound
		default:
			return nil, err
		}
	}
	return &post, nil
}

func (s *PostStore) Update(ctx context.Context, post *postsEntity.Post) error {
	query := `
		UPDATE posts
		SET title = $1, content = $2, version = version + 1 
		WHERE id = $3 AND version = $4
		RETURNING version
	`

	ctx, cancel := context.WithTimeout(ctx, storage.QueryTimeoutDuration)
	defer cancel()

	err := s.Db.QueryRowContext(
		ctx,
		query,
		post.Title,
		post.Content,
		post.ID,
		post.Version,
	).Scan(&post.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return storage.ErrNotFound
		default:
			return err
		}
	}
	return nil
}

func (s *PostStore) Delete(ctx context.Context, postID int64) error {
	query := `
		DELETE FROM posts
		WHERE posts.id = $1
	`

	ctx, cancel := context.WithTimeout(ctx, storage.QueryTimeoutDuration)
	defer cancel()

	res, err := s.Db.ExecContext(
		ctx,
		query,
		postID,
	)

	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return storage.ErrNotFound
	}

	return nil
}
