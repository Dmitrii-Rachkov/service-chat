package db

import (
	"database/sql"

	"service-chat/internal/db/entity"
)

// Здесь интерфейсы для слоя базы данных нашего приложения

// Authorization - интерфейс авторизации
type Authorization interface {
	CreateUser(user entity.User) (int, error)
	GetUser(user entity.User) (*entity.User, error)
}

// Chat - интерфейс для чатов
type Chat interface {
	CreateChat(in entity.ChatAdd) (int, error)
	GetChat(in entity.ChatGet) ([]entity.Chat, error)
	DeleteChat(in entity.ChatDelete) ([]entity.DeletedChats, error)
}

// Message - интерфейс для сообщений
type Message interface {
	AddMessage(in entity.MessageAdd) (int, error)
	UpdateMessage(in entity.MessageUpdate) (int, error)
	GetMessage(in entity.MessageGet) ([]entity.Message, error)
	DeleteMessage(in entity.MessageDel) ([]entity.DelMsg, error)
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
		Chat:          NewChatsPostgres(db),
		Message:       NewMessagePostgres(db),
	}
}
