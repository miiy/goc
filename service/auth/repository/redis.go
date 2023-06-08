package repository

import (
	"context"
	"github.com/miiy/goc/redis"
	"time"
)

type redisRepository struct {
	rdb redis.UniversalClient
	AuthTokenRepository
}

func NewRedisRepository(rdb redis.UniversalClient) AuthTokenRepository {
	return &redisRepository{
		rdb: rdb,
	}
}

func (r *redisRepository) Get(ctx context.Context, key string) (string, error) {
	return r.rdb.Get(ctx, key).Result()
}

func (r *redisRepository) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.rdb.Set(ctx, key, value, expiration).Err()
}

func (r *redisRepository) Del(ctx context.Context, key string) error {
	return r.rdb.Del(ctx, key).Err()
}
