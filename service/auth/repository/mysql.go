package repository

import (
	"context"
	"fmt"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/service/auth/entity"
)

type mysqlAuthRepository struct {
	db *gorm.DB
	AuthRepository
}

func NewAuthRepository(db *gorm.DB) AuthRepository {
	return &mysqlAuthRepository{
		db: db,
	}
}

func (r *mysqlAuthRepository) Create(ctx context.Context, user *entity.User) error {
	err := r.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *mysqlAuthRepository) Update(ctx context.Context, id uint64, v *entity.User, columns ...string) (rowsAffected int64, err error) {
	result := r.db.WithContext(ctx).Select(columns).Where("id = ?", id).Updates(v)
	return result.RowsAffected, result.Error
}

func (r *mysqlAuthRepository) First(ctx context.Context, id uint64, columns ...string) (*entity.User, error) {
	var user entity.User
	err := r.db.WithContext(ctx).Select(columns).First(&user, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &user, nil
}

func (r *mysqlAuthRepository) FirstByUsername(ctx context.Context, username string, columns ...string) (*entity.User, error) {
	var item entity.User
	err := r.db.WithContext(ctx).Select(columns).Where("username=?", username).First(&item).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &item, nil
}

func (r *mysqlAuthRepository) UserExist(ctx context.Context, column, value string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Where(fmt.Sprintf("%s=?", column), value).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}
