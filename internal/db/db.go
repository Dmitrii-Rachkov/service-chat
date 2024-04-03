package db

import "database/sql"

// Здесь интерфейсы для слоя базы данных нашего приложения

// Authorization - интерфейс авторизации
type Authorization interface {
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
	return &DB{}
}
