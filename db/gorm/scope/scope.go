package scope

import (
	"gorm.io/gorm"
)

func ScopeOfUser(userId []int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(userId) == 1 {
			db.Where("user_id = ?", userId)
		}
		return db.Where("user_id IN (?)", userId)
	}
}

func ScopeOfStatus(status []int64) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if len(status) == 1 {
			db.Where("status = ?", status)
		}
		return db.Where("status IN (?)", status)
	}
}
