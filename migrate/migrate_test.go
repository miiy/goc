package migrate

import (
	"strings"
	"testing"
	"testing/fstest"
)

func TestValidate(t *testing.T) {
	valid := Config{
		Database: DatabaseConfig{
			Driver:   "mysql",
			Host:     "localhost",
			Port:     "3306",
			Database: "test",
		},
		MigrationsTable: "test_migrations",
		Source: fstest.MapFS{
			"000001_initial.up.sql": &fstest.MapFile{Data: []byte("SELECT 1;")},
		},
		SourcePath: ".",
	}

	if err := validate(valid, " UP "); err != nil {
		t.Fatalf("validate() error = %v", err)
	}
	if err := validate(valid, "force"); err == nil || !strings.Contains(err.Error(), "unsupported command") {
		t.Fatalf("validate() error = %v, want unsupported command", err)
	}
}

func TestRunRejectsUnsupportedDriver(t *testing.T) {
	config := Config{
		Database: DatabaseConfig{
			Driver:   "postgres",
			Host:     "localhost",
			Port:     "5432",
			Database: "test",
		},
		MigrationsTable: "test_migrations",
		Source: fstest.MapFS{
			"000001_initial.up.sql": &fstest.MapFile{Data: []byte("SELECT 1;")},
		},
		SourcePath: ".",
	}

	if err := Run(config, "up"); err == nil || !strings.Contains(err.Error(), "not supported") {
		t.Fatalf("Run() error = %v, want unsupported driver", err)
	}
}
