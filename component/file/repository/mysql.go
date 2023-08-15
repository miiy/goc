package repository

import (
	"context"
	"github.com/miiy/goc/component/file/entity"
	"github.com/miiy/goc/db/gorm"
)

type mysqlFileRepository struct {
	db *gorm.DB
}

func NewFileRepository(db *gorm.DB) FileRepository {
	return &mysqlFileRepository{
		db: db,
	}
}

func ScopeOfFileSys(sys int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("sys = ?", sys)
	}
}

func ScopeOfFileCat(catId int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("cat_id = ?", catId)
	}
}

func (r *mysqlFileRepository) CreateFile(ctx context.Context, i *entity.File) (*entity.File, error) {
	err := r.db.WithContext(ctx).Create(&i).Error
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (r *mysqlFileRepository) UpdateFile(ctx context.Context, id int64, i *entity.File, selects interface{}) (int64, error) {
	ret := r.db.WithContext(ctx).Model(&entity.File{}).Where("id = ?", id).Select(selects).Updates(i)
	err := ret.Error
	if err != nil {
		return 0, err
	}
	return ret.RowsAffected, nil
}

func (r *mysqlFileRepository) DeleteFileById(ctx context.Context, id int64) error {
	err := r.db.WithContext(ctx).Delete(&entity.File{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *mysqlFileRepository) GetFileById(ctx context.Context, id int64, scopes ...func(*gorm.DB) *gorm.DB) (*entity.File, error) {
	var i entity.File
	err := r.db.WithContext(ctx).Model(&entity.File{}).Scopes(scopes...).First(&i, id).Error
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *mysqlFileRepository) FindCount(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var count int64
	r.db.WithContext(ctx).Model(&entity.File{}).Scopes(scopes...).Count(&count)
	return count, nil
}

func (r *mysqlFileRepository) Find(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) ([]*entity.File, error) {
	var items []*entity.File
	r.db.WithContext(ctx).Model(&entity.File{}).Scopes(scopes...).Find(&items)

	return items, nil
}
