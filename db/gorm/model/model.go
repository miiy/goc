package model

import (
	"database/sql"
	"fmt"
	"gorm.io/gorm/schema"
	"strings"
	"sync"
	"time"
)

type Model struct {
	Id              int64        `gorm:"column:id"`
	CreateTime      time.Time    `gorm:"column:create_time;autoCreateTime"`
	UpdateTime      time.Time    `gorm:"column:update_time;autoUpdateTime"`
	DeleteTime      sql.NullTime `gorm:"column:delete_time"`
}

const (
	FieldNameFormatWithQuote = "`%s`"
	FieldNameFormatWithPlaceHolder = "`%s` = ?"
)

var FieldNameExpectAutoSet = []string{
	"id",
	"create_time",
	"update_time",
	"delete_time",
}

func FieldDBNames(dest interface{}, excepts []string) ([]string, error) {
	s, err := schema.Parse(dest, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		return nil, err
	}
	var dbNames []string

	for _, field := range s.Fields {
		var find bool
		for _, except := range excepts {
			if find = except == field.DBName; find{
				break
			}
		}
		if find {
			continue
		}
		dbNames = append(dbNames, field.DBName)
	}

	return dbNames, nil
}

func FieldNameFormat(fields []string, format string) string {
	var r []string
	for _, field := range fields {
		r = append(r, fmt.Sprintf(format, field))
	}
	return strings.Join(r, ", ")
}
