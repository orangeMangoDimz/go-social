package cache

import (
	"context"

	usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"
)

func NewMockStore() Storage {
	return Storage{
		Users: &MockUserStore{},
	}
}

type MockUserStore struct{}

func (m *MockUserStore) Get(ctx context.Context, userID int64) (*usersEntity.User, error) {
	return &usersEntity.User{}, nil
}

func (m *MockUserStore) Set(ctx context.Context, user *usersEntity.User) error {
	return nil
}
