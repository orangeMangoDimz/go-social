package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/orangeMangoDimz/go-social/store"
	"github.com/redis/go-redis/v9"
)

type UserStore struct {
	rdb *redis.Client
}

const UserExpTime = time.Minute

func (s *UserStore) Get(ctx context.Context, userID int64) (*store.User, error) {
	// Return nil if Redis client is not available
	if s.rdb == nil {
		return nil, nil
	}

	cacheKey := fmt.Sprintf("user-%v", userID)

	data, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == redis.Nil {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	var user store.User
	if data != "" {
		if err := json.Unmarshal([]byte(data), &user); err != nil {
			return nil, err
		}
	}

	return &user, nil

}

func (s *UserStore) Set(ctx context.Context, user *store.User) error {
	// Return nil if Redis client is not available (no-op)
	if s.rdb == nil {
		return nil
	}

	cacheKey := fmt.Sprintf("user-%v", user.ID)

	json, err := json.Marshal(user)
	if err != nil {
		return err
	}

	return s.rdb.SetEx(ctx, cacheKey, json, UserExpTime).Err()
}
