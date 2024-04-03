package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"service-chat/internal/config"
)

const driver = "postgres"

// NewPostgresDB - соединение с базой данных postgres
func NewPostgresDB(cfg *config.Config) (*sql.DB, error) {
	// Параметры подключения к базе данных
	url := fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=%v",
		cfg.Database.Username,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.SSLMode)

	db, err := sql.Open(driver, url)
	if err != nil {
		return nil, err
	}

	// Проверяем, что можем достучаться до базы данных
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Устанавливаем максимальное количество одновременных подключений к базе
	db.SetMaxOpenConns(cfg.Database.Connections)

	return db, nil
}
