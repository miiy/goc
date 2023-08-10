package repository

import (
	"context"
	"github.com/miiy/goc/auth"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthenticateRepository interface {
	RetrieveByIdentifier(ctx context.Context, identifier string, value interface{}) (*auth.AuthenticatedUser, error)
	//RetrieveByCredentials(ctx context.Context, identifier string, value interface{}, password string) (*auth.AuthenticatedUser, error)
}

type authenticateRepository struct {
	db  *gorm.DB
	rdb redis.UniversalClient
}

func (r *authRepository) RetrieveByIdentifier(ctx context.Context, identifier string, value interface{}) (*auth.AuthenticatedUser, error) {
	var user auth.User
	err := r.db.WithContext(ctx).Where("?=?", identifier, value).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &auth.AuthenticatedUser{
		ID:       user.ID,
		Username: user.Username,
	}, nil
}
