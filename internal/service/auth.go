package service

import (
	"service-chat/internal/db"
	"service-chat/internal/encryption"
	"service-chat/internal/entity"
)

type AuthService struct {
	repo db.Authorization
}

// NewAuthService - конструктор для работы со слоем сервиса
func NewAuthService(repo db.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// CreateUser - реализуем интерфейс функции CreateUser передав данные со слоя бизнес логики на слой базы данных
func (s *AuthService) CreateUser(user entity.User) (int, error) {
	// Шифруем пароль, чтобы хранить в базе не в открытом виде
	passwordHash, err := encryption.Encrypt(user.Password)
	if err != nil {
		return 0, err
	}
	user.Password = passwordHash

	// Передаём в слой базы данных
	return s.repo.CreateUser(user)
}
