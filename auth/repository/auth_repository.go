package repository

import (
	"context"
	"fmt"
	"github.com/miiy/goc/auth"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type AuthRepository interface {
	Create(ctx context.Context, user *auth.User) error
	Update(ctx context.Context, id uint64, user *auth.User, columns ...string) (rowsAffected int64, err error)
	First(ctx context.Context, id uint64, columns ...string) (*auth.User, error)
	FirstByUsername(ctx context.Context, username string, columns ...string) (*auth.User, error)
	FirstByMpOpenid(ctx context.Context, openid string, columns ...string) (*auth.User, error)
	UserExist(ctx context.Context, column, value string) (bool, error)
}

type authRepository struct {
	db  *gorm.DB
	rdb redis.UniversalClient
	AuthRepository
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &authRepository{
		db: db,
	}
}

func (r *authRepository) Create(ctx context.Context, user *auth.User) error {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *authRepository) Update(ctx context.Context, id uint64, v *auth.User, columns ...string) (rowsAffected int64, err error) {
	result := r.db.WithContext(ctx).Select(columns).Where("id = ?", id).Updates(v)
	return result.RowsAffected, result.Error
}

func (r *authRepository) First(ctx context.Context, id uint64, columns ...string) (*auth.User, error) {
	var user auth.User
	err := r.db.WithContext(ctx).Select(columns).First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *authRepository) FirstByUsername(ctx context.Context, username string, columns ...string) (*auth.User, error) {
	var item auth.User
	err := r.db.WithContext(ctx).Select(columns).Where("username=?", username).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *authRepository) FirstByMpOpenid(ctx context.Context, openid string, columns ...string) (*auth.User, error) {
	var item auth.User
	err := r.db.WithContext(ctx).Select(columns).Where("mp_openid=?", openid).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *authRepository) UserExist(ctx context.Context, column, value string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&auth.User{}).Where(fmt.Sprintf("%s=?", column), value).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
