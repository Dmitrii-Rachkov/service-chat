package db

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"service-chat/internal/config"
)

// MigrateUp - накатываем миграцию базы данных
func MigrateUp(cfg *config.Config, log *slog.Logger) error {
	const op = "db.MigrateUp"
	databaseURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode)
	m, err := migrate.New("file://./internal/db/schema", databaseURL)
	if err != nil {
		return fmt.Errorf("error path: %s, error: %w", op, err)
	}

	if errUp := m.Up(); errors.Is(errUp, migrate.ErrNoChange) {
		log.Info("MigrateDB", slog.String("migrateUp", "No changes"))
	} else if errUp != nil && !errors.Is(errUp, migrate.ErrNoChange) {
		return fmt.Errorf("error path: %s, errorUp: %w", op, errUp)
	}

	return nil
}

// MigrateDown - откатываем миграцию базы данных
func MigrateDown(cfg *config.Config) error {
	const op = "db.MigrateDown"
	databaseURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode)
	m, err := migrate.New("file://./internal/db/schema", databaseURL)
	if err != nil {
		return fmt.Errorf("error path: %s, error: %w", op, err)
	}
	if errDown := m.Down(); errDown != nil {
		return fmt.Errorf("error path: %s, error: %w", op, errDown)
	}

	return nil
}
