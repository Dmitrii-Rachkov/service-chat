package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"service-chat/internal/config"
)

// MigrateUp - накатываем миграцию базы данных
func MigrateUp(cfg *config.Config) error {
	databaseURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode)
	m, err := migrate.New("file://./internal/db/schema", databaseURL)
	if err != nil {
		return err
	}
	if errUp := m.Up(); err != nil {
		return errUp
	}

	return nil
}

// MigrateDown - откатываем миграцию базы данных
func MigrateDown(cfg *config.Config) error {
	databaseURL := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode)
	m, err := migrate.New("file://./internal/db/schema", databaseURL)
	if err != nil {
		return err
	}
	if errDown := m.Down(); err != nil {
		return errDown
	}

	return nil
}
