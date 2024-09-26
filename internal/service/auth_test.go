package service

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"service-chat/internal/db"
	"service-chat/internal/db/entity"
	mockRepo "service-chat/internal/db/mocks"
	"service-chat/internal/dto"
)

func TestAuthService_CreateUser(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockRepo.MockAuthorization, dataDB entity.User)

	tests := []struct {
		name    string
		user    dto.SignUpRequest
		dataDB  entity.User
		mock    mockBehaviour
		wantID  int
		wantErr error
		encrypt bool
	}{
		{
			name: "Success",
			user: dto.SignUpRequest{
				Username: "Andrey",
				Password: "adgui*",
			},
			dataDB: entity.User{
				Username: "Andrey",
				Password: "CiRA9gEG",
			},
			mock: func(s *mockRepo.MockAuthorization, dataDB entity.User) {
				s.EXPECT().CreateUser(dataDB).Return(1, nil)
			},
			wantID:  1,
			wantErr: nil,
		},
		{
			name: "Error encrypt",
			user: dto.SignUpRequest{
				Username: "Andrey",
				Password: "adgui*",
			},
			mock:    func(s *mockRepo.MockAuthorization, dataDB entity.User) {},
			wantID:  0,
			wantErr: errors.New("encryption.Encrypt: crypto/aes: invalid key size 0"),
			encrypt: true,
		},
		{
			name: "Other error",
			user: dto.SignUpRequest{
				Username: "Andrey",
				Password: "adgui*",
			},
			dataDB: entity.User{
				Username: "Andrey",
				Password: "CiRA9gEG",
			},
			mock: func(s *mockRepo.MockAuthorization, dataDB entity.User) {
				s.EXPECT().CreateUser(dataDB).Return(0, errors.New("some error"))
			},
			wantID:  0,
			wantErr: errors.New("some error"),
		},
		{
			name: "Empty request",
			user: dto.SignUpRequest{
				Username: "",
				Password: "",
			},
			mock:    func(s *mockRepo.MockAuthorization, dataDB entity.User) {},
			wantID:  0,
			wantErr: errors.New("username or password is empty"),
		},
		{
			name:    "Nil request",
			user:    dto.SignUpRequest{},
			mock:    func(s *mockRepo.MockAuthorization, dataDB entity.User) {},
			wantID:  0,
			wantErr: errors.New("username or password is empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Инициализируем зависимости

			//Инициализируем контролер для мока сервиса
			ctrl := gomock.NewController(t)
			// Завершаем работу контролера после выполнения каждого теста
			defer ctrl.Finish()

			// Создаём моки базы данных авторизации
			mockAuth := mockRepo.NewMockAuthorization(ctrl)

			// Передаём структуру пользователя
			tt.mock(mockAuth, tt.dataDB)

			// Создаём объект базы данных в который передадим наш мок авторизации
			repository := &db.DB{Authorization: mockAuth}

			// Создаём экземпляр сервиса авторизации
			serviceAuth := NewAuthService(repository)

			// Проверяем ожидаемый и актуальный результат
			if tt.encrypt {
				acID, acErr := serviceAuth.CreateUser(tt.user)
				assert.Equal(t, tt.wantID, acID)
				assert.Equal(t, tt.wantErr.Error(), acErr.Error())
			} else {
				t.Setenv("MY_SECRET", "abc&1*~#^2^#s0^=)^^7%b34")
				acID, acErr := serviceAuth.CreateUser(tt.user)
				assert.Equal(t, tt.wantID, acID)
				assert.Equal(t, tt.wantErr, acErr)
			}
		})
	}
}

func TestAuthService_GenerateToken(t *testing.T) {
	type mockBehaviour func(s *mockRepo.MockAuthorization, dataDB entity.User)

	tests := []struct {
		name    string
		user    dto.SignInRequest
		dataDB  entity.User
		mock    mockBehaviour
		wantErr error
		decrypt bool
	}{
		{
			name: "Success",
			user: dto.SignInRequest{
				Username: "Andrey",
				Password: "adgui*",
			},
			dataDB: entity.User{
				Username: "Andrey",
				Password: "adgui*",
			},
			mock: func(s *mockRepo.MockAuthorization, dataDB entity.User) {
				s.EXPECT().GetUser(dataDB).Return(&entity.User{
					Username: "Andrey",
					Password: "CiRA9gEG",
				}, nil)
			},
		},
		{
			name: "Error GetUser from db",
			user: dto.SignInRequest{
				Username: "Andrey",
				Password: "adgui*",
			},
			dataDB: entity.User{
				Username: "Andrey",
				Password: "adgui*",
			},
			mock: func(s *mockRepo.MockAuthorization, dataDB entity.User) {
				s.EXPECT().GetUser(dataDB).Return(nil, errors.New("some error"))
			},
			wantErr: errors.New("some error"),
		},
		{
			name: "Empty request",
			user: dto.SignInRequest{
				Username: "",
				Password: "",
			},
			mock:    func(s *mockRepo.MockAuthorization, dataDB entity.User) {},
			wantErr: errors.New("username or password is empty"),
		},
		{
			name:    "Nil request",
			mock:    func(s *mockRepo.MockAuthorization, dataDB entity.User) {},
			wantErr: errors.New("username or password is empty"),
		},
		{
			name: "Error decrypt password",
			user: dto.SignInRequest{
				Username: "Andrey",
				Password: "adgui*",
			},
			dataDB: entity.User{
				Username: "Andrey",
				Password: "adgui*",
			},
			mock: func(s *mockRepo.MockAuthorization, dataDB entity.User) {
				s.EXPECT().GetUser(dataDB).Return(&entity.User{
					Username: "Andrey",
					Password: "CiRA9gEG",
				}, nil)
			},
			wantErr: errors.New("failed to decrypt password: encryption.Decrypt: crypto/aes: invalid key size 0"),
			decrypt: true,
		},
		{
			name: "incorrect password",
			user: dto.SignInRequest{
				Username: "Andrey",
				Password: "adgui*",
			},
			dataDB: entity.User{
				Username: "Andrey",
				Password: "adgui*",
			},
			mock: func(s *mockRepo.MockAuthorization, dataDB entity.User) {
				s.EXPECT().GetUser(dataDB).Return(&entity.User{
					Username: "Andrey",
					Password: "Fail",
				}, nil)
			},
			wantErr: errors.New("incorrect password"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Инициализируем контролер для мока сервиса
			ctrl := gomock.NewController(t)
			// Завершаем работу контролера после выполнения каждого теста
			defer ctrl.Finish()

			// Создаём моки базы данных авторизации
			mockAuth := mockRepo.NewMockAuthorization(ctrl)

			// Передаём структуру пользователя
			tt.mock(mockAuth, tt.dataDB)

			// Создаём объект базы данных в который передадим наш мок авторизации
			repository := &db.DB{Authorization: mockAuth}

			// Создаём экземпляр сервиса авторизации
			serviceAuth := NewAuthService(repository)

			// Проверяем ожидаемый и актуальный результат
			if tt.decrypt {
				acToken, acErr := serviceAuth.GenerateToken(tt.user)
				assert.Empty(t, acToken)
				assert.Equal(t, tt.wantErr.Error(), acErr.Error())
			} else if tt.wantErr != nil {
				t.Setenv("MY_SECRET", "abc&1*~#^2^#s0^=)^^7%b34")
				acToken, acErr := serviceAuth.GenerateToken(tt.user)
				assert.Empty(t, acToken)
				assert.Equal(t, tt.wantErr, acErr)
			} else {
				t.Setenv("MY_SECRET", "abc&1*~#^2^#s0^=)^^7%b34")
				acToken, acErr := serviceAuth.GenerateToken(tt.user)
				assert.NotEmpty(t, acToken)
				assert.Equal(t, tt.wantErr, acErr)
			}
		})
	}
}

func TestAuthService_ParseToken(t *testing.T) {
	// Создаём подписанный jwt token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
		1,
	})
	jwtToken, errToken := token.SignedString([]byte(signingKey))
	assert.NoError(t, errToken)

	// Невалидный Claims токен
	tokenErrClaims := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MjY5NDI0NzQsImlhdCI6MTcyNjk0MjE3NCwidXNlcl9pZCI6MX0.eus0W6Rj87AkltBluogH-zTthlHovVjljQOnPcz82z0"

	// Невалидный SignMethod
	tokenRSA, errRSA := createTokenRSA()
	assert.NoError(t, errRSA)

	tests := []struct {
		name    string
		token   string
		wantID  int
		wantErr error
	}{
		{
			name:    "Success",
			token:   jwtToken,
			wantID:  1,
			wantErr: nil,
		},
		{
			name:    "Error claims",
			token:   tokenErrClaims,
			wantID:  0,
			wantErr: errors.New("token has invalid claims: token is expired"),
		},
		{
			name:    "Error SignMethod",
			token:   tokenRSA,
			wantID:  0,
			wantErr: errors.New("token is unverifiable: error while executing keyfunc: unexpected signing method"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//Инициализируем контролер для мока сервиса
			ctrl := gomock.NewController(t)
			// Завершаем работу контролера после выполнения каждого теста
			defer ctrl.Finish()

			// Создаём моки базы данных авторизации
			mockAuth := mockRepo.NewMockAuthorization(ctrl)

			// Создаём объект базы данных в который передадим наш мок авторизации
			repository := &db.DB{Authorization: mockAuth}

			// Создаём экземпляр сервиса авторизации
			serviceAuth := NewAuthService(repository)

			// Проверяем ожидаемый и актуальный результат
			acID, acErr := serviceAuth.ParseToken(tt.token)
			assert.Equal(t, tt.wantID, acID)
			if acErr != nil {
				assert.Equal(t, tt.wantErr.Error(), acErr.Error())
			} else {
				assert.NoError(t, acErr)
			}
		})
	}
}

type JWT struct {
	privateKey []byte
	publicKey  []byte
}

func NewJWT(privateKey []byte, publicKey []byte) JWT {
	return JWT{
		privateKey: privateKey,
		publicKey:  publicKey,
	}
}

func (j JWT) createRSA(ttl time.Duration, content interface{}) (string, error) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM(j.privateKey)
	if err != nil {
		return "", fmt.Errorf("createRSA: parse key: %w", err)
	}

	now := time.Now().UTC()

	claims := make(jwt.MapClaims)
	claims["dat"] = content             // Our custom data.
	claims["exp"] = now.Add(ttl).Unix() // The expiration time after which the token must be disregarded.
	claims["iat"] = now.Unix()          // The time at which the token was issued.
	claims["nbf"] = now.Unix()          // The time before which the token must be disregarded.

	token, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
	if err != nil {
		return "", fmt.Errorf("create: sign token: %w", err)
	}

	return token, nil
}

//go:embed mocks/id_rsa
var prvKey []byte

//go:embed mocks/id_rsa.pub
var pubKey []byte

func createTokenRSA() (token string, err error) {
	jwtToken := NewJWT(prvKey, pubKey)

	// Create a new JWT token.
	tok, err := jwtToken.createRSA(time.Hour, "Can be anything")
	if err != nil {
		log.Fatalln(err)
	}
	return tok, nil
}
