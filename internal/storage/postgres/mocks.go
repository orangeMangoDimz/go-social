package postgres

import (
	"context"
	"database/sql"
	"time"

	usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"
	"github.com/orangeMangoDimz/go-social/internal/storage"
)

func NewMockStore() storage.Storage {
	return storage.Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct {
}

func (m *MockUserStore) Create(ctx context.Context, tx *sql.Tx, user *usersEntity.User) error {
	return nil
}

func (m *MockUserStore) GetById(ctx context.Context, userID int64) (*usersEntity.User, error) {
	return &usersEntity.User{}, nil
}

func (m *MockUserStore) GetByEmail(ctx context.Context, userEmail string) (*usersEntity.User, error) {
	return &usersEntity.User{}, nil
}

func (m *MockUserStore) CreateAndInvite(ctx context.Context, user *usersEntity.User, token string, invitationExp time.Duration) error {
	return nil
}

func (m *MockUserStore) Activate(ctx context.Context, token string) error {
	return nil
}

func (m *MockUserStore) Delete(ctx context.Context, userID int64) error {
	return nil
}
