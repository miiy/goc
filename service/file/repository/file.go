package repository

import (
	"context"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/db/gorm/scope"
)

type FileRepositoryImpl struct {
	db *gorm.DB
}

type FileRepository interface {
	CreateFile(ctx context.Context, file *File) (*File, error)
	UpdateFile(ctx context.Context, id int64, file *File, selects interface{}) (int64, error)
	DeleteFileById(ctx context.Context, id int64) error
	GetFileById(ctx context.Context, id int64) (*File, error)
	FindCount(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error)
	Find(ctx context.Context, scopes ...func(db *gorm.DB) *gorm.DB) ([]*File, error)
}

type File struct {
	gorm.Model
	Sys      int    `gorm:"column:sys"`
	Type     int    `gorm:"column:type" json:"type"`
	UserId   int64  `gorm:"column:user_id"`
	FileType int    `gorm:"column:file_type"`
	Name     string `gorm:"column:name"`
	Ext      string `gorm:"column:ext"`
	Path     string `gorm:"column:path"`
	Hash     string `gorm:"column:hash"`
	Status   int    `gorm:"column:status"`
}

const (
	StatusDefault = 0
	StatusActive  = 1
	StatusDisable = 2
)

const (
	RecommendNo  = 0
	RecommendYes = 1
)

var (
	FieldNames                string
	FieldNamesExpectAutoSet   string
	FieldNamesWithPlaceHolder string
)

func NewFileRepository(db *gorm.DB) FileRepository {
	return &FileRepositoryImpl{
		db: db,
	}
}

func ScopeOfFileUser(userId int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("user_id = ?", userId)
	}
}

func ScopeOfFileCategory(categoryId int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("category_id = ?", categoryId)
	}
}

func (r *FileRepositoryImpl) CreateFile(ctx context.Context, i *File) (*File, error) {
	err := r.db.WithContext(ctx).Create(&i).Error
	if err != nil {
		return nil, err
	}
	return i, nil
}

func (r *FileRepositoryImpl) UpdateFile(ctx context.Context, id int64, i *File, selects interface{}) (int64, error) {
	ret := r.db.WithContext(ctx).Model(&File{}).Where("id = ?", id).Select(selects).Updates(i)
	err := ret.Error
	if err != nil {
		return 0, err
	}
	return ret.RowsAffected, nil
}

func (r *FileRepositoryImpl) DeleteFileById(ctx context.Context, id int64) error {
	err := r.db.WithContext(ctx).Delete(&File{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *FileRepositoryImpl) GetFileById(ctx context.Context, id int64) (*File, error) {
	var i File
	err := r.db.WithContext(ctx).Model(&File{}).Scopes(scope.ScopeOfStatus([]int64{StatusActive})).First(&i, id).Error
	if err != nil {
		return nil, err
	}
	return &i, nil
}

func (r *FileRepositoryImpl) FindCount(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) (int64, error) {
	var count int64
	r.db.WithContext(ctx).Model(&File{}).Scopes(scopes...).Count(&count)
	return count, nil
}

func (r *FileRepositoryImpl) Find(ctx context.Context, scopes ...func(*gorm.DB) *gorm.DB) ([]*File, error) {
	var items []*File
	r.db.WithContext(ctx).Model(&File{}).Scopes(scopes...).Find(&items)

	return items, nil
}
