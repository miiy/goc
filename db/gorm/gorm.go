package gorm

import (
	"gorm.io/gorm"
	"time"
)

type DB = gorm.DB

type Model struct {
	ID        int64          `gorm:"column:id;primarykey"`
	CreatedAt time.Time      `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time      `gorm:"column:updated_at;autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;index"`
}

var ErrRecordNotFound = gorm.ErrRecordNotFound
