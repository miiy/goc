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

type GormDB = gorm.DB

type Database struct {
	DB   *sql.DB
	Gorm *gorm.DB
}

type Config struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type Options struct {
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	MaxOpenConns    int
}

type Option func(*Options)

var defaultOption = Options{
	ConnMaxLifetime: time.Minute * 3,
	MaxIdleConns:    10,
	MaxOpenConns:    100,
}

var (
	ErrRecordNotFound = gorm.ErrRecordNotFound
	ErrCreateError    = errors.New("create error")
	ErrUpdateError    = errors.New("update error")
)

func WithConnMaxLifetime(t time.Duration) Option {
	return func(c *Options) {
		c.ConnMaxLifetime = t
	}
}

func WithMaxIdleConns(n int) Option {
	return func(c *Options) {
		c.MaxIdleConns = n
	}
}

func WithMaxOpenConns(n int) Option {
	return func(c *Options) {
		c.MaxOpenConns = n
	}
}

// NewDatabase
// dns refer https://github.com/go-sql-driver/mysql for details
func NewDatabase(c Config, opts ...Option) (*Database, error) {
	if c.Driver == "mysql" {
		dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true&loc=Local&multiStatements=true", c.Username, c.Password, c.Host, c.Port, c.Database)
		db, err := sql.Open(c.Driver, dsn)
		if err != nil {
			return nil, err
		}

		c := defaultOption
		for _, o := range opts {
			o(&c)
		}

		db.SetConnMaxLifetime(c.ConnMaxLifetime)
		db.SetMaxIdleConns(c.MaxIdleConns)
		db.SetMaxOpenConns(c.MaxOpenConns)

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

		return &Database{
			DB:   db,
			Gorm: gormDB,
		}, nil
	}

	return nil, errors.New("database: driver not support")
}
