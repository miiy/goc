package pagination

import (
	"math"
)

type Pagination struct {
	Total       int64
	PerPage     int64
	CurrentPage int64
	LastPage    int64
	From        int64
	To          int64
}

var (
	PerPageDefault int64 = 20
	MaxPerPage     int64 = 100
)

func NewPagination(page, perPage, total int64) *Pagination {
	if perPage <= 0 {
		perPage = PerPageDefault
	}

	if perPage > MaxPerPage {
		perPage = MaxPerPage
	}

	lastPage := int64(math.Ceil(float64(total) / float64(perPage)))
	if page <= 0 {
		page = 1
	}
	// total=0, lastPage = 0
	if lastPage <= 0 {
		lastPage = 1
	}
	if page > lastPage {
		page = lastPage
	}

	var offset = (page - 1) * perPage

	return &Pagination{
		Total:       total,
		PerPage:     perPage,
		CurrentPage: page,
		LastPage:    lastPage,
		From:        offset,
		To:          offset + perPage,
	}
}
