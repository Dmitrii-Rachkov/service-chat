package handler

import (
	"bytes"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"service-chat/internal/dto"
	"service-chat/internal/service"
	mockService "service-chat/internal/service/mocks"
)

// Юнит тест для обработчика Регистрация пользователя SignUp
func TestHandler_SignUp(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockService.MockAuthorization, user dto.SignUpRequest)

	// Тестовая таблица с данными
	testTable := []struct {
		name                 string
		inputBody            string
		inputUser            dto.SignUpRequest
		mockBehavior         mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"username": "Andrey", "password": "adgui*"}`,
			inputUser: dto.SignUpRequest{
				Username: "Andrey",
				Password: "adgui*",
			},
			// Реализуем поведение мока, возвращаем userID и nil ошибку
			mockBehavior: func(s *mockService.MockAuthorization, user dto.SignUpRequest) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Create user is successful, id: 1"}`,
		},
		{
			name:      "Required field username is missing",
			inputBody: `{"password": "adgui*"}`,
			// Реализуем поведение мока, возвращаем userID и nil ошибку
			mockBehavior:         func(s *mockService.MockAuthorization, user dto.SignUpRequest) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Username is a required field"}`,
		},
		{
			name:      "Required field password is missing",
			inputBody: `{"username": "Andrey"}`,
			// Реализуем поведение мока, возвращаем userID и nil ошибку
			mockBehavior:         func(s *mockService.MockAuthorization, user dto.SignUpRequest) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Password is a required field"}`,
		},
		{
			name:      "All required field is missing",
			inputBody: `{}`,
			// Реализуем поведение мока, возвращаем userID и nil ошибку
			mockBehavior:         func(s *mockService.MockAuthorization, user dto.SignUpRequest) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Username is a required field, Field Password is a required field"}`,
		},
		{
			name: "Request body is nil",
			// Реализуем поведение мока, возвращаем userID и nil ошибку
			mockBehavior:         func(s *mockService.MockAuthorization, user dto.SignUpRequest) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Empty request"}`,
		},
		{
			name:      "Error decode",
			inputBody: `{"username": , "password": "adgui*"}`,
			// Реализуем поведение мока, возвращаем userID и nil ошибку
			mockBehavior:         func(s *mockService.MockAuthorization, user dto.SignUpRequest) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:      "Username max 20 chars",
			inputBody: `{"username": "Andreytuoplkjhgsdtyw", "password": "adgui*"}`,
			inputUser: dto.SignUpRequest{
				Username: "Andreytuoplkjhgsdtyw",
				Password: "adgui*",
			},
			// Реализуем поведение мока, возвращаем userID и nil ошибку
			mockBehavior: func(s *mockService.MockAuthorization, user dto.SignUpRequest) {
				s.EXPECT().CreateUser(user).Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Create user is successful, id: 1"}`,
		},
		{
			name:      "Username 21 chars",
			inputBody: `{"username": "Andreytuoplkjhgsdtywk", "password": "adgui*"}`,
			// Реализуем поведение мока, возвращаем userID и nil ошибку
			mockBehavior:         func(s *mockService.MockAuthorization, user dto.SignUpRequest) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Username cannot exceed 20 characters"}`,
		},
		{
			name:      "Username forbidden characters",
			inputBody: `{"username": "Andrey@", "password": "adgui*"}`,
			// Реализуем поведение мока, возвращаем userID и nil ошибку
			mockBehavior:         func(s *mockService.MockAuthorization, user dto.SignUpRequest) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Username must not contain symbols !@#$\u0026*()?"}`,
		},
		{
			name:      "Password min 6 chars",
			inputBody: `{"username": "Andrey", "password": "adgui"}`,
			// Реализуем поведение мока, возвращаем userID и nil ошибку
			mockBehavior:         func(s *mockService.MockAuthorization, user dto.SignUpRequest) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Password must contain at least 6 characters"}`,
		},
		{
			name:      "Password max 12 chars",
			inputBody: `{"username": "Andrey", "password": "qwertyuiopasd"}`,
			// Реализуем поведение мока, возвращаем userID и nil ошибку
			mockBehavior:         func(s *mockService.MockAuthorization, user dto.SignUpRequest) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Password cannot exceed 12 characters"}`,
		},
		{
			name:      "Password contains special character",
			inputBody: `{"username": "Andrey", "password": "qwerty"}`,
			// Реализуем поведение мока, возвращаем userID и nil ошибку
			mockBehavior:         func(s *mockService.MockAuthorization, user dto.SignUpRequest) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Password must contain Latin letters and Arabic numerals, as well as the symbols @#$\u0026*()"}`,
		},
		{
			name:      "Unique violation username",
			inputBody: `{"username": "Andrey", "password": "adgui*"}`,
			inputUser: dto.SignUpRequest{
				Username: "Andrey",
				Password: "adgui*",
			},
			// Реализуем поведение мока, возвращаем userID и nil ошибку
			mockBehavior: func(s *mockService.MockAuthorization, user dto.SignUpRequest) {
				s.EXPECT().CreateUser(user).Return(0, errors.New("unique_violation"))
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"User already exists"}`,
		},
		{
			name:      "Other error",
			inputBody: `{"username": "Andrey", "password": "adgui*"}`,
			inputUser: dto.SignUpRequest{
				Username: "Andrey",
				Password: "adgui*",
			},
			// Реализуем поведение мока, возвращаем userID и nil ошибку
			mockBehavior: func(s *mockService.MockAuthorization, user dto.SignUpRequest) {
				s.EXPECT().CreateUser(user).Return(0, errors.New("fail"))
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Failed to create user"}`,
		},
	}

	// Итерируемся по нашей тестовой таблице
	for _, tt := range testTable {
		// Запускаем тесты параллельно
		t.Run(tt.name, func(t *testing.T) {
			// Инициализируем зависимости

			//Инициализируем контролер для мока сервиса
			ctrl := gomock.NewController(t)
			// Завершаем работу контролера после выполнения каждого теста
			defer ctrl.Finish()

			// Создаём моки сервиса авторизации
			mockAuth := mockService.NewMockAuthorization(ctrl)

			// Передаём структуру пользователя
			tt.mockBehavior(mockAuth, tt.inputUser)

			// Создаём объект сервиса в который передадим наш мок авторизации
			services := &service.Service{Authorization: mockAuth}

			// Создаём экземпляр обработчика
			handler := NewHandler(services)

			// Мокируем логгер
			mockLog := slog.New(slog.NewJSONHandler(io.Discard, nil))

			// Инициализируем сервер

			// Инициализируем тестовый endPoint по которому будет вызываться тестовый обработчик
			r := chi.NewRouter()
			r.Post("/sign-up", handler.SignUp(mockLog))

			// Готовим тестовый запрос
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/sign-up", bytes.NewBufferString(tt.inputBody))

			// Выполняем запрос
			r.ServeHTTP(w, req)

			// Сравниваем ожидаемый и актуальный результат
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
