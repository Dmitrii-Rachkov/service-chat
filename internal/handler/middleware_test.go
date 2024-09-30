package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"service-chat/internal/service"
	mockService "service-chat/internal/service/mocks"
)

func TestHandler_AuthMiddleware(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockService.MockAuthorization, token string)

	// Инициализируем зависимости

	//Инициализируем контролер для мока сервиса
	ctrl := gomock.NewController(t)
	// Завершаем работу контролера после выполнения каждого теста
	defer ctrl.Finish()

	// Создаём моки сервиса авторизации
	mockAuth := mockService.NewMockAuthorization(ctrl)

	// Создаём объект сервиса в который передадим наш мок авторизации
	services := &service.Service{Authorization: mockAuth}

	// Создаём экземпляр обработчика
	handler := NewHandler(services)

	// Инициализируем сервер

	// Инициализируем тестовый endPoint, в котором используем AuthMiddleware
	// В этом endPoint из контекста будем забирать userID и записывать в ответ
	r := chi.NewRouter()
	r.Use(handler.AuthMiddleware)

	r.Post("/protected", func(w http.ResponseWriter, r *http.Request) {
		idCtx, _ := GetUserID(r.Context())
		render.JSON(w, r, OK(strconv.Itoa(idCtx)))
	})

	testTable := []struct {
		name                 string
		headerName           string
		headerValue          string
		token                string
		mockBehavior         mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:        "OK",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			// Реализуем поведение мока, даём на вход token и возвращаем userID и nil ошибку
			mockBehavior: func(s *mockService.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(1, nil)
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"OK","message":"1"}`,
		},
		{
			name:                 "Invalid Header Name",
			headerName:           "",
			headerValue:          "Bearer token",
			token:                "token",
			mockBehavior:         func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"Error","error":"Authorization header is empty"}`,
		},
		{
			name:                 "Invalid Header Value",
			headerName:           "Authorization",
			headerValue:          "Bear token",
			token:                "token",
			mockBehavior:         func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"Error","error":"Invalid authorization header"}`,
		},
		{
			name:                 "Empty Token",
			headerName:           "Authorization",
			headerValue:          "Bearer ",
			token:                "token",
			mockBehavior:         func(s *mockService.MockAuthorization, token string) {},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"Error","error":"Token is empty"}`,
		},
		{
			name:        "Parse Error",
			headerName:  "Authorization",
			headerValue: "Bearer token",
			token:       "token",
			mockBehavior: func(s *mockService.MockAuthorization, token string) {
				s.EXPECT().ParseToken(token).Return(0, errors.New("parse error"))
			},
			expectedStatusCode:   200,
			expectedResponseBody: `{"status":"Error","error":"parse error"}`,
		},
	}

	// Итерируемся по нашей тестовой таблице
	for _, tt := range testTable {
		// Запускаем тесты параллельно
		t.Run(tt.name, func(t *testing.T) {
			// Передаём токен пользователя
			tt.mockBehavior(mockAuth, tt.token)

			// Формируем запрос у которого в headers есть токен авторизации
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/protected", nil)
			req.Header.Set(tt.headerName, tt.headerValue)

			// Выполняем запрос
			r.ServeHTTP(w, req)

			// Сравниваем ожидаемый и актуальный результат
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func Test_GetUserID(t *testing.T) {
	// Функция для создания тестового контекста
	var getContext = func(id int) context.Context {
		ctx := context.WithValue(context.Background(), userCtx, id)
		return ctx
	}

	testTable := []struct {
		name          string
		ctx           context.Context
		expectedID    int
		isFail        bool
		expectedError error
	}{
		{
			name:       "OK",
			ctx:        getContext(1),
			expectedID: 1,
		},
		{
			name:          "Not found",
			ctx:           context.Background(),
			isFail:        true,
			expectedError: errors.New("user id not found"),
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			id, err := GetUserID(tt.ctx)
			if tt.isFail {
				assert.Error(t, err)
				assert.Equal(t, tt.expectedError, err)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, tt.expectedID, id)
		})
	}
}
