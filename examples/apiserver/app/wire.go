//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/miiy/goc/auth/jwt"
	"github.com/miiy/goc/db"
	"github.com/miiy/goc/examples/apiserver/config"
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc/redis"
	authRepo "github.com/miiy/goc/service/auth/repository"
	authServer "github.com/miiy/goc/service/auth/server"
	"gorm.io/gorm"
)

func InitApp(conf string) (*App, func(), error) {
	panic(wire.Build(
		config.NewConfig,
		logger.NewLogger, providerLoggerOption,
		db.NewDatabase, providerDatabase, providerDatabaseOption, providerGorm,
		redis.NewRedis, providerRedisOptions,
		jwt.NewJWTAuth, providerJwtAuthOptions,
		wire.NewSet(authRepo.NewAuthRepository, authServer.NewAuthServiceServer, authRepo.NewRedisRepository),
		NewApp,
	))
}

func providerGorm(db *db.Database) *gorm.DB {
	return db.Gorm
}
func providerRedisOptions(config *config.Config) *redis.Options {
	return &redis.Options{
		Addrs:    config.Redis.Addrs,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	}
}

func providerJwtAuthOptions(config *config.Config) *jwt.Options {
	return &jwt.Options{
		Secret:    config.Jwt.Secret,
		Issuer:    config.Jwt.Issuer,
		ExpiresIn: config.Jwt.ExpiresIn,
	}
}

func providerLoggerOption(config *config.Config) []logger.Option {
	return []logger.Option{
		//logger.WithEnv(environment.Environment(config.App.Env)),
	}
}

func providerDatabase(config *config.Config) db.Config {
	return db.Config{
		Driver:   config.Database.Driver,
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		Username: config.Database.Username,
		Password: config.Database.Password,
		Database: config.Database.Database,
	}
}

func providerDatabaseOption(config *config.Config) []db.Option {
	return []db.Option{
		//db.WithEnv(environment.Environment(config.App.Env)),
	}
}
