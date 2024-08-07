package service

import (
	"service-chat/internal/db"
	"service-chat/internal/dto"
)

// Здесь интерфейсы для слоя бизнес-логики нашего приложения

// Authorization - интерфейс авторизации
type Authorization interface {
	// CreateUser - функция для создания нового пользователя в базе и вернуть его id или ошибку
	CreateUser(user dto.SignUpRequest) (int, error)
	// GenerateToken - создаём токен для авторизации пользователя
	GenerateToken(user dto.SignInRequest) (string, error)
	// ParseToken - анализируем jwt token
	ParseToken(token string) (int, error)
}

// Chat - интерфейс для чатов
type Chat interface {
	// CreateChat - создаём чат между пользователями
	CreateChat(in dto.ChatAdd) (int, error)
}

// Message - интерфейс для сообщений
type Message interface {
	// AddMessage - отправить сообщение в чат от лица пользователя
	AddMessage(in dto.MessageAdd) (int, error)
	UpdateMessage(in dto.MessageUpdate) (int, error)
}

// Service - собирает все наши интерфейсы в одном месте
type Service struct {
	Authorization
	Chat
	Message
}

// NewService - конструктор сервиса
func NewService(db *db.DB) *Service {
	return &Service{
		Authorization: NewAuthService(db.Authorization),
		Chat:          NewChatService(db.Chat),
		Message:       NewMessageService(db.Message),
	}
}
