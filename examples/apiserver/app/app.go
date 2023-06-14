package app

import (
	"github.com/miiy/goc/auth/jwt"
	"github.com/miiy/goc/db"
	"github.com/miiy/goc/examples/apiserver/config"
	"github.com/miiy/goc/logger"
	"github.com/miiy/goc/redis"
	authpb "github.com/miiy/goc/service/auth/api/v1"
)

type App struct {
	Config     *config.Config
	Database   *db.Database
	Redis      redis.UniversalClient
	Logger     logger.Logger
	JwtAuth    *jwt.JWTAuth
	AuthServer authpb.AuthServiceServer
}

var app *App

func NewApp(c *config.Config, db *db.Database, rdb redis.UniversalClient, l logger.Logger, j *jwt.JWTAuth, as authpb.AuthServiceServer) *App {
	app = &App{
		Config:     c,
		Database:   db,
		Redis:      rdb,
		Logger:     l,
		JwtAuth:    j,
		AuthServer: as,
	}
	return app
}

func Instance() *App {
	return app
}
