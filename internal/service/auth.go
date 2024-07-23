package service

import (
	"service-chat/internal/db"
	"service-chat/internal/db/entity"
	"service-chat/internal/dto"
	"service-chat/internal/encryption"
)

type AuthService struct {
	repo db.Authorization
}

// NewAuthService - конструктор для работы со слоем сервиса
func NewAuthService(repo db.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

// CreateUser - реализуем интерфейс функции CreateUser передав данные со слоя бизнес логики на слой базы данных
func (s *AuthService) CreateUser(user dto.SignUpRequest) (int, error) {
	// Шифруем пароль, чтобы хранить в базе не в открытом виде
	passwordHash, err := encryption.Encrypt(user.Password)
	if err != nil {
		return 0, err
	}

	dataDB := entity.User{
		Username: user.Username,
		Password: passwordHash,
	}

	// Передаём в слой базы данных
	return s.repo.CreateUser(dataDB)
}
