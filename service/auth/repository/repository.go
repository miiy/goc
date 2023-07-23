package repository

import (
	"context"
	"github.com/miiy/goc/service/auth/entity"
	"time"
)

type AuthRepository interface {
	Create(ctx context.Context, user *entity.User) error
	Update(ctx context.Context, id uint64, user *entity.User, columns ...string) (rowsAffected int64, err error)
	First(ctx context.Context, id uint64, columns ...string) (*entity.User, error)
	FirstByUsername(ctx context.Context, username string, columns ...string) (*entity.User, error)
	FirstByMpOpenid(ctx context.Context, openid string, columns ...string) (*entity.User, error)
	UserExist(ctx context.Context, column, value string) (bool, error)
}

type AuthTokenRepository interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
}
