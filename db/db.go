package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type DB struct {
	opts options
	db   *sql.DB
	gorm *gorm.DB
}

type Config struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// NewDB
// dns refer https://github.com/go-sql-driver/mysql for details
func NewDB(c Config, opts ...Option) (*DB, error) {
	if c.Driver == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&multiStatements=true", c.Username, c.Password, c.Host, c.Port, c.Database)
		db, err := sql.Open(c.Driver, dsn)
		if err != nil {
			return nil, err
		}

		dopts := defaultOption
		for _, o := range opts {
			o.apply(&dopts)
		}

		db.SetConnMaxLifetime(dopts.connMaxLifetime)
		db.SetMaxIdleConns(dopts.maxIdleConns)
		db.SetMaxOpenConns(dopts.maxOpenConns)

		newLogger := logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容——译者注）
			logger.Config{
				SlowThreshold:             time.Second, // 慢 SQL 阈值
				LogLevel:                  logger.Info, // 日志级别
				IgnoreRecordNotFoundError: true,        // 忽略ErrRecordNotFound（记录未找到）错误
				Colorful:                  false,       // 禁用彩色打印
			},
		)

		// gorm
		gormDB, err := gorm.Open(mysql.New(mysql.Config{
			Conn: db,
		}), &gorm.Config{
			Logger: newLogger,
			//NamingStrategy: schema.NamingStrategy{
			//	TablePrefix:   "",
			//	SingularTable: true,
			//	NameReplacer:  nil,
			//	NoLowerCase:   false,
			//},
		})
		if err != nil {
			return nil, err
		}

		return &DB{
			opts: dopts,
			db:   db,
			gorm: gormDB,
		}, nil
	}

	return nil, errors.New("database: driver not support")
}

func (db *DB) DB() *sql.DB {
	return db.db
}

func (db *DB) Gorm() *gorm.DB {
	return db.gorm
}
