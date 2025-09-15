package cache

import (
	"context"

	usersEntity "github.com/orangeMangoDimz/go-social/internal/entities/users"
	usersCache "github.com/orangeMangoDimz/go-social/internal/storage/cache/users"
	"github.com/redis/go-redis/v9"
)

type Storage struct {
	Users interface {
		Get(context.Context, int64) (*usersEntity.User, error)
		Set(context.Context, *usersEntity.User) error
	}
}

func NewRedisStorage(rdb *redis.Client) Storage {
	return Storage{
		Users: &usersCache.UserStore{Rdb: rdb},
	}
}
