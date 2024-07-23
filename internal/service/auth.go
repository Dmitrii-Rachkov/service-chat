package service

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"service-chat/internal/db"
	"service-chat/internal/db/entity"
	"service-chat/internal/dto"
	"service-chat/internal/encryption"
)

const (
	tokenTTL   = 5 * time.Second
	signingKey = "qWeRtYuIoP123456789#@&*"
)

// Расширяем стандартный токен
type tokenClaims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"user_id"`
}

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

// GenerateToken - реализуем создание jwt token
func (s *AuthService) GenerateToken(user dto.SignInRequest) (string, error) {
	// Получаем пользователя из базы данных
	dataDB := entity.User{
		Username: user.Username,
		Password: user.Password,
	}

	userDB, err := s.repo.GetUser(dataDB)
	if err != nil {
		return "", err
	}

	// Проверяем, что пользователь ввёл верный пароль
	passwordDecode, errDecode := encryption.Decrypt(userDB.Password)
	if errDecode != nil {
		return "", fmt.Errorf("failed to decrypt password: %w", errDecode)
	}
	if passwordDecode != user.Password {
		return "", fmt.Errorf("incorrect password")
	}

	// Если пользователь существует, генерируем token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		userDB.Id,
	})

	// Создаём подписанный jwt token
	jwtToken, errToken := token.SignedString([]byte(signingKey))
	if errToken != nil {
		return "", fmt.Errorf("failed to generation jwt token")
	}

	return jwtToken, nil
}
