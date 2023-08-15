package repository

import (
	"context"
	"github.com/miiy/goc/component/file/entity"
	"gorm.io/gorm"
)

type FileRepository interface {
	CreateFile(ctx context.Context, file *entity.File) (*entity.File, error)
	UpdateFile(ctx context.Context, id int64, file *entity.File, selects interface{}) (int64, error)
	DeleteFileById(ctx context.Context, id int64) error
	GetFileById(ctx context.Context, id int64, scopes ...func(*gorm.DB) *gorm.DB) (*entity.File, error)
	FindCount(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error)
	Find(ctx context.Context, scopes ...func(db *gorm.DB) *gorm.DB) ([]*entity.File, error)
}
