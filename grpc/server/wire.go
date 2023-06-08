//go:build wireinject
// +build wireinject

package main

//import (
//	"docxlib.com/pkg/database"
//	"docxlib.com/pkg/environment"
//	"docxlib.com/pkg/jwtauth"
//	"docxlib.com/pkg/log"
//	"docxlib.com/service/user/internal/app"
//	"docxlib.com/service/user/internal/config"
//	"github.com/google/wire"
//)
//
//func newApp(conf string) (*app.App, func(), error) {
//	panic(
//		wire.Build(
//			app.ProviderSet,
//			config.ProviderSet,
//			log.ProviderSet, providerLoggerOption,
//			providerDatabase, providerDatabaseOption, database.NewDatabase,
//			jwtauth.ProviderSet,
//		),
//	)
//}
//
//func providerLoggerOption() []log.Option {
//	return []log.Option{
//		log.WithEnv(environment.DEVELOPMENT),
//	}
//}
//
//func providerDatabase(config *config.Config) database.Config {
//	return database.Config{
//		Driver:   config.Database.Driver,
//		Host:     config.Database.Host,
//		Port:     config.Database.Port,
//		Username: config.Database.Username,
//		Password: config.Database.Password,
//		Database: config.Database.Database,
//	}
//}
//
//func providerDatabaseOption(config *config.Config) []database.Option {
//	return []database.Option{
//		database.WithEnv(environment.Environment(config.App.Env)),
//	}
//}
