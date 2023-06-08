package gorm

import (
	"database/sql"
	"gorm.io/gorm"
	"time"
)

type DB = gorm.DB

type Model struct {
	Id         int64        `gorm:"column:id;primarykey"`
	CreateTime time.Time    `gorm:"column:create_time;autoCreateTime"`
	UpdateTime time.Time    `gorm:"column:update_time;autoUpdateTime"`
	DeleteTime sql.NullTime `gorm:"column:delete_time;index"`
}

var ErrRecordNotFound = gorm.ErrRecordNotFound
