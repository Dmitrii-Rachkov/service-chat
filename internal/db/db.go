package db

import (
	"database/sql"

	"service-chat/internal/db/entity"
)

// Здесь интерфейсы для слоя базы данных нашего приложения

// Authorization - интерфейс авторизации
type Authorization interface {
	CreateUser(user entity.User) (int, error)
}

// Chat - интерфейс для чатов
type Chat interface {
}

// Message - интерфейс для сообщений
type Message interface {
}

// DB - собирает все наши интерфейсы в одном месте
type DB struct {
	Authorization
	Chat
	Message
}

// NewDB - конструктор базы данных
func NewDB(db *sql.DB) *DB {
	return &DB{
		Authorization: NewAuthPostgres(db),
	}
}
