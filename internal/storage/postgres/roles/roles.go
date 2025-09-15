package roles

import (
	"context"
	"database/sql"
	"errors"

	usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"
	"github.com/orangeMangoDimz/go-social/internal/storage"
)

type RoleStore struct {
	Db *sql.DB
}

func (r *RoleStore) GetByName(ctx context.Context, roleName string) (*usersEntity.Role, error) {
	query := `
		SELECT id, name, level, description
		FROM roles
		WHERE name = $1
	`

	ctx, cancel := context.WithTimeout(ctx, storage.QueryTimeoutDuration)
	defer cancel()

	var role usersEntity.Role
	err := r.Db.QueryRowContext(
		ctx,
		query,
		roleName,
	).Scan(
		&role.ID,
		&role.Name,
		&role.Level,
		&role.Description,
	)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, storage.ErrNotFound
		default:
			return nil, err
		}
	}
	return &role, nil
}
