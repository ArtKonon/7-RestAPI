package repository

import (
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	_ = db
	return nil
}

func RunMigrations(databaseURL, migrationsPath string) error {
	if databaseURL == "" {
		databaseURL = os.Getenv("DATABASE_URL")
	}
	if databaseURL == "" {
		return fmt.Errorf("database URL is empty")
	}
	if migrationsPath == "" {
		migrationsPath = os.Getenv("MIGRATIONS_PATH")
	}
	if migrationsPath == "" {
		migrationsPath = "migrations"
	}

	m, err := migrate.New("file://"+migrationsPath, databaseURL)
	if err != nil {
		return fmt.Errorf("create migrate runner: %w", err)
	}
	defer func() {
		_, _ = m.Close()
	}()

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return fmt.Errorf("apply migrations: %w", err)
	}
	return nil
}
