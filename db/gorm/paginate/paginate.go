package paginate

import (
	"gorm.io/gorm"
	"math"
)

var (
	DefaultPageSize    = 20
	DefaultMaxPageSize = 100
)

func Paginate(page, pageSize, maxPageSize, total int) (scope func(db *gorm.DB) *gorm.DB, totalPage int) {
	if page == 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = DefaultPageSize
	}
	if maxPageSize <= 0 {
		maxPageSize = DefaultMaxPageSize
	}
	if pageSize > maxPageSize {
		pageSize = maxPageSize
	}

	offset := (page - 1) * pageSize
	totalPage = int(math.Ceil(float64(total) / float64(pageSize)))
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(offset).Limit(pageSize)
	}, totalPage
}
