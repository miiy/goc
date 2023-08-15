package entity

import "github.com/miiy/goc/db/gorm"

var (
	FieldNames                string
	FieldNamesExpectAutoSet   string
	FieldNamesWithPlaceHolder string
)

type File struct {
	gorm.Model
	SysId    int64
	CatId    int64
	ItemId   int64
	UserId   int64
	FileType int
	Name     string
	Ext      string
	Path     string
	Hash     string
	Status   int
}

const (
	FileStatusDefault = 0
	FileStatusActive  = 1
	FileStatusDisable = 2
)
