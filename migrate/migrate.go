package migrate

import (
	"errors"
	"fmt"
	"io/fs"
	"strings"
)

// DatabaseConfig contains the database connection settings used by a migration.
type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

// Config configures a Migrator using versioned migrations from Source.
type Config struct {
	Database        DatabaseConfig
	MigrationsTable string
	Source          fs.FS
	SourcePath      string
}

// Run applies migrations using the configured database driver.
func Run(config Config, command string) error {
	if err := validate(config, command); err != nil {
		return err
	}

	switch strings.ToLower(strings.TrimSpace(config.Database.Driver)) {
	case "mysql":
		return runMySQL(config, strings.ToLower(strings.TrimSpace(command)))
	default:
		return fmt.Errorf("migrate: database driver %q is not supported", config.Database.Driver)
	}
}

func validate(config Config, command string) error {
	if config.Source == nil {
		return errors.New("migrate: source filesystem is required")
	}
	if strings.TrimSpace(config.Database.Database) == "" {
		return errors.New("migrate: database name is required")
	}
	if strings.TrimSpace(config.Database.Host) == "" {
		return errors.New("migrate: database host is required")
	}
	if strings.TrimSpace(config.Database.Port) == "" {
		return errors.New("migrate: database port is required")
	}
	if strings.TrimSpace(config.MigrationsTable) == "" {
		return errors.New("migrate: migrations table is required")
	}
	switch strings.ToLower(strings.TrimSpace(command)) {
	case "up", "down":
		return nil
	default:
		return fmt.Errorf("migrate: unsupported command %q (expected up or down)", command)
	}
}
