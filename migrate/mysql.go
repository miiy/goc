package migrate

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/go-sql-driver/mysql"
	gomigrate "github.com/golang-migrate/migrate/v4"
	migratemysql "github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

func runMySQL(config Config, command string) error {
	driverConfig := mysql.Config{
		User:                 config.Database.Username,
		Passwd:               config.Database.Password,
		Net:                  "tcp",
		Addr:                 config.Database.Host + ":" + config.Database.Port,
		DBName:               config.Database.Database,
		ParseTime:            true,
		Loc:                  time.Local,
		MultiStatements:      true,
		AllowNativePasswords: true,
	}
	database, err := sql.Open("mysql", driverConfig.FormatDSN())
	if err != nil {
		return fmt.Errorf("migrate: open mysql database: %w", err)
	}
	if err := database.Ping(); err != nil {
		_ = database.Close()
		return fmt.Errorf("migrate: ping mysql database: %w", err)
	}

	sourceDriver, err := iofs.New(config.Source, config.SourcePath)
	if err != nil {
		_ = database.Close()
		return fmt.Errorf("migrate: open source: %w", err)
	}
	databaseDriver, err := migratemysql.WithInstance(database, &migratemysql.Config{
		DatabaseName:    config.Database.Database,
		MigrationsTable: config.MigrationsTable,
	})
	if err != nil {
		_ = sourceDriver.Close()
		_ = database.Close()
		return fmt.Errorf("migrate: create mysql driver: %w", err)
	}
	engine, err := gomigrate.NewWithInstance("iofs", sourceDriver, config.Database.Database, databaseDriver)
	if err != nil {
		_ = sourceDriver.Close()
		_ = databaseDriver.Close()
		return fmt.Errorf("migrate: create migrator: %w", err)
	}

	var runErr error
	switch command {
	case "up":
		runErr = engine.Up()
	case "down":
		runErr = engine.Down()
	}
	if errors.Is(runErr, gomigrate.ErrNoChange) {
		runErr = nil
	}

	sourceErr, databaseErr := engine.Close()
	return errors.Join(
		runErr,
		wrapCloseError("source", sourceErr),
		wrapCloseError("database driver", databaseErr),
	)
}

func wrapCloseError(resource string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("migrate: close %s: %w", resource, err)
}
