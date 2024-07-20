package service

import (
	"service-chat/internal/db"
	"service-chat/internal/entity"
)

// Здесь интерфейсы для слоя бизнес-логики нашего приложения

// Authorization - интерфейс авторизации
type Authorization interface {
	// CreateUser - функция для создания нового пользователя в базе и вернуть его id или ошибку
	CreateUser(user entity.User) (int, error)
}

// Chat - интерфейс для чатов
type Chat interface {
}

// Message - интерфейс для сообщений
type Message interface {
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
	}
}
