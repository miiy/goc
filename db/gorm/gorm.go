package gorm

import (
	"time"

	"gorm.io/gorm"
)

type DB = gorm.DB
type Config = gorm.Config
type Dialector = gorm.Dialector
type Option = gorm.Option

type Model struct {
	ID        int64          `gorm:"column:id;primarykey"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

var ErrRecordNotFound = gorm.ErrRecordNotFound

func Open(dialector gorm.Dialector, opts ...gorm.Option) (*DB, error) {
	return gorm.Open(dialector, opts...)
}
