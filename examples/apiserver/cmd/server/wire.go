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
	"github.com/miiy/goc/examples/apiserver/app"
	"github.com/miiy/goc/examples/apiserver/config"
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func initApp(conf string) (*app.App, func(), error) {
	panic(wire.Build(
		config.NewConfig,
		wire.NewSet(logger.NewLogger, providerLoggerOption, providerZap),
		wire.NewSet(db.NewDB, providerDBConfig, providerDBOption, providerGorm),
		wire.NewSet(redis.NewRedis, providerRedisOptions),
		wire.NewSet(jwt.NewJWTAuth, providerJwtAuthOptions),
		wire.NewSet(authRepo.NewAuthRepository, authServer.NewAuthServiceServer, authRepo.NewTokenRepository, providerMiniProgram, providerUser),
		app.NewApp,
	))
}

func providerLoggerOption() []logger.Option {
	return nil
}

func providerZap(logger logger.Logger) *zap.Logger {
	return logger.ZapLogger()
}

func providerDBConfig(config *config.Config) db.Config {
	return db.Config{
		Driver:   config.Database.Driver,
		Host:     config.Database.Host,
		Port:     config.Database.Port,
		Username: config.Database.Username,
		Password: config.Database.Password,
		Database: config.Database.Database,
	}
}

func providerDBOption(config *config.Config) []db.Option {
	return nil
}

func providerRedisOptions(config *config.Config) *redis.Options {
	return &redis.Options{
		Addrs:    config.Redis.Addrs,
		Password: config.Redis.Password,
		DB:       config.Redis.DB,
	}
}

func providerGorm(db *db.DB) *gorm.DB {
	return db.Gorm()
}

func providerJwtAuthOptions(config *config.Config) *jwt.Options {
	return &jwt.Options{
		Secret:    config.Jwt.Secret,
		Issuer:    config.Jwt.Issuer,
		ExpiresIn: config.Jwt.ExpiresIn,
	}
}

func providerMiniProgram(config *config.Config) (*miniprogram.MiniProgram, error) {
	return nil, nil
}

func providerUser(authRepository authRepo.AuthRepository) auth.UserProvider {
	return authRepository
}
