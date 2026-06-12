package mysql

import (
	gocgorm "github.com/miiy/goc/db/gorm"
	"gorm.io/driver/mysql"
)

type Config = mysql.Config
type Dialector = mysql.Dialector

func New(config Config) gocgorm.Dialector {
	return mysql.New(config)
}

func Open(dsn string) gocgorm.Dialector {
	return mysql.Open(dsn)
}
