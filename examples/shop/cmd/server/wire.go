//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/miiy/goc/auth"
	"github.com/miiy/goc/auth/jwt"
	authRepo "github.com/miiy/goc/component/auth/repository"
	authServer "github.com/miiy/goc/component/auth/server"
	"github.com/miiy/goc/contrib/sdk/wechat/miniprogram"
	"github.com/miiy/goc/db"
	"github.com/miiy/goc/db/gorm"
	"github.com/miiy/goc/examples/shop/internal/app"
	"github.com/miiy/goc/examples/shop/internal/config"
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc/redis"
	"go.uber.org/zap"
)

func initApp(conf string) (*app.App, func(), error) {
	panic(wire.Build(
		config.NewConfig,
		wire.NewSet(logger.NewLogger, provideLoggerOption, provideZap),
		wire.NewSet(db.NewDB, provideDBConfig, provideDBOption, provideGorm),
		wire.NewSet(redis.NewRedis, provideRedisOptions),
		wire.NewSet(jwt.NewJWTAuth, provideJwtAuthOptions),
		wire.NewSet(authRepo.NewAuthRepository, authServer.NewAuthServiceServer, authRepo.NewTokenRepository, provideMiniProgram),
		wire.Bind(new(auth.UserProvider), new(authRepo.AuthRepository)),
		app.NewApp,
	))
}

func provideLoggerOption() []logger.Option {
	return nil
}

func provideZap(logger logger.Logger) *zap.Logger {
	return logger.ZapLogger()
}

func provideDBConfig(config *config.Config) db.Config {
	return db.Config{
		Driver:   config.Database.Driver,
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		Username: config.Database.Username,
		Password: config.Database.Password,
		Database: config.Database.Database,
	}
}

func provideDBOption(config *config.Config) []db.Option {
	return nil
}

func provideRedisOptions(config *config.Config) *redis.Options {
	return &redis.Options{
		Addrs:    config.Redis.Addrs,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	}
}

func provideGorm(db *db.DB) *gorm.DB {
	return db.Gorm()
}

func provideJwtAuthOptions(config *config.Config) *jwt.Options {
	return &jwt.Options{
		Secret:    config.Jwt.Secret,
		Issuer:    config.Jwt.Issuer,
		ExpiresIn: config.Jwt.ExpiresIn,
	}
}

func provideMiniProgram(config *config.Config) (*miniprogram.MiniProgram, error) {
	return nil, nil
}
