package handler

import (
	"context"
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

	"service-chat/internal/db/entity"
	"service-chat/internal/dto"
	"service-chat/internal/service"
	mockService "service-chat/internal/service/mocks"
)

func TestHandler_ChatAdd(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockService.MockChat, chat dto.ChatAdd)

	// Тестовая таблица с данными
	testTable := []struct {
		name                 string
		inputBody            string
		inputChat            dto.ChatAdd
		mockBehavior         mockBehaviour
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"chat_name": "chat_1","users": [1,2]}`,
			inputChat: dto.ChatAdd{
				ChatName: "chat_1",
				Users:    []int64{1, 2},
			},
			mockBehavior: func(s *mockService.MockChat, chat dto.ChatAdd) {
				s.EXPECT().CreateChat(chat).Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Chat created successfully, id: 1"}`,
		},
		{
			name:                 "Required field chat_name is missing",
			inputBody:            `{"users": [1,2]}`,
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field ChatName is a required field"}`,
		},
		{
			name:                 "Required field chat_name is missing",
			inputBody:            `{"chat_name": "chat_1"}`,
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Users is a required field"}`,
		},
		{
			name:                 "All required field is missing",
			inputBody:            `{}`,
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field ChatName is a required field, Field Users is a required field"}`,
		},
		{
			name:                 "Request body is nil",
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Empty request"}`,
		},
		{
			name:                 "Error decode",
			inputBody:            `{"chat_name": ,"users": [1,2]}`,
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:      "Chat_name max 20 chars",
			inputBody: `{"chat_name": "qwertyuiopasdfghjklz","users": [1,2]}`,
			inputChat: dto.ChatAdd{
				ChatName: "qwertyuiopasdfghjklz",
				Users:    []int64{1, 2},
			},
			mockBehavior: func(s *mockService.MockChat, chat dto.ChatAdd) {
				s.EXPECT().CreateChat(chat).Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Chat created successfully, id: 1"}`,
		},
		{
			name:                 "Chat_name max 21 chars",
			inputBody:            `{"chat_name": "qwertyuiopasdfghjklzx","users": [1,2]}`,
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field ChatName cannot exceed 20 characters"}`,
		},
		{
			name:                 "Chat_name forbidden characters",
			inputBody:            `{"chat_name": "chat_1#","users": [1,2]}`,
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field ChatName must not contain symbols !@#$\u0026*()?"}`,
		},
		{
			name:                 "Chat_name min 6 chars",
			inputBody:            `{"chat_name": "chat","users": [1,2]}`,
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field ChatName must contain at least 6 characters"}`,
		},
		{
			name:                 "Empty users",
			inputBody:            `{"chat_name": "chat_1","users": []}`,
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Users must have at least 2 elements"}`,
		},
		{
			name:                 "One user",
			inputBody:            `{"chat_name": "chat_1","users": [1]}`,
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Users must have at least 2 elements"}`,
		},
		{
			name:                 "Chat_name is nil",
			inputBody:            `{"chat_name": "","users": [1,2]}`,
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field ChatName is a required field"}`,
		},
		{
			name:      "Unique violation chat_name",
			inputBody: `{"chat_name": "chat_1","users": [1,2]}`,
			inputChat: dto.ChatAdd{
				ChatName: "chat_1",
				Users:    []int64{1, 2},
			},
			mockBehavior: func(s *mockService.MockChat, chat dto.ChatAdd) {
				s.EXPECT().CreateChat(chat).Return(1, errors.New("unique_violation"))
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Chat already exists"}`,
		},
		{
			name:      "Other error",
			inputBody: `{"chat_name": "chat_1","users": [1,2]}`,
			inputChat: dto.ChatAdd{
				ChatName: "chat_1",
				Users:    []int64{1, 2},
			},
			mockBehavior: func(s *mockService.MockChat, chat dto.ChatAdd) {
				s.EXPECT().CreateChat(chat).Return(1, errors.New("example error"))
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Failed to create chat: example error"}`,
		},
	}

	// Итерируемся по нашей тестовой таблице
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			// Инициализируем зависимости

			//Инициализируем контролер для мока сервиса
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Создаём моки сервиса чат
			mockChat := mockService.NewMockChat(ctrl)

			// Передаём структуру пользователя
			tt.mockBehavior(mockChat, tt.inputChat)

			// Создаём объект сервиса в который передадим наш мок авторизации
			services := &service.Service{Chat: mockChat}

			// Создаём экземпляр обработчика
			handler := NewHandler(services)

			// Мокируем логгер
			mockLog := slog.New(slog.NewJSONHandler(io.Discard, nil))

			// Инициализируем сервер

			// Инициализируем тестовый endPoint по которому будет вызываться тестовый обработчик
			r := chi.NewRouter()
			r.Post("/chats/add", handler.ChatAdd(mockLog))

			// Готовим тестовый запрос
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/chats/add", strings.NewReader(tt.inputBody))

			// Выполняем запрос
			r.ServeHTTP(w, req)

			// Сравниваем ожидаемый и актуальный результат
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}

}

func TestHandler_ChatDelete(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehavior func(s *mockService.MockChat, chat dto.ChatDelete)

	// Тестовая таблица с данными
	testTable := []struct {
		name                 string
		inputBody            string
		inputChat            dto.ChatDelete
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK one chat",
			inputBody: `{"chat_ids": [1]}`,
			inputChat: dto.ChatDelete{
				ChatIds: &[]int64{1},
			},
			mockBehavior: func(s *mockService.MockChat, chat dto.ChatDelete) {
				s.EXPECT().DeleteChat(chat, 1).Return([]entity.DeletedChats{
					{
						ChatID: 1,
						Result: "Chat successfully deleted",
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Result of deleted chats","del_chats_list":[{"chat_id":1,"result":"Chat successfully deleted"}]}`,
		},
		{
			name:      "OK many chats",
			inputBody: `{"chat_ids": [1,2,3]}`,
			inputChat: dto.ChatDelete{
				ChatIds: &[]int64{1, 2, 3},
			},
			mockBehavior: func(s *mockService.MockChat, chat dto.ChatDelete) {
				s.EXPECT().DeleteChat(chat, 1).Return([]entity.DeletedChats{
					{
						ChatID: 1,
						Result: "Chat successfully deleted",
					},
					{
						ChatID: 2,
						Result: "Chat successfully deleted",
					},
					{
						ChatID: 3,
						Result: "Chat successfully deleted",
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Result of deleted chats","del_chats_list":[{"chat_id":1,"result":"Chat successfully deleted"},{"chat_id":2,"result":"Chat successfully deleted"},{"chat_id":3,"result":"Chat successfully deleted"}]}`,
		},
		{
			name:                 "Required field chat_ids is missing",
			inputBody:            `{}`,
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatDelete) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field ChatIds is a required field"}`,
		},
		{
			name:                 "Empty chat_ids",
			inputBody:            `{"chat_ids": []}`,
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatDelete) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field ChatIds must contain at least 1 characters"}`,
		},
		{
			name:                 "Request body is nil",
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatDelete) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Empty request"}`,
		},
		{
			name:                 "Error decode",
			inputBody:            `{"chat_ids":}`,
			mockBehavior:         func(s *mockService.MockChat, chat dto.ChatDelete) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:      "Other error",
			inputBody: `{"chat_ids": [1]}`,
			inputChat: dto.ChatDelete{
				ChatIds: &[]int64{1},
			},
			mockBehavior: func(s *mockService.MockChat, chat dto.ChatDelete) {
				s.EXPECT().DeleteChat(chat, 1).Return(nil, errors.New("some error"))
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Failed to delete chats: some error"}`,
		},
	}

	// Итерируемся по нашей тестовой таблице
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			// Инициализируем зависимости

			//Инициализируем контролер для мока сервиса
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Создаём моки сервиса чат
			mockChat := mockService.NewMockChat(ctrl)

			// Передаём структуру пользователя
			tt.mockBehavior(mockChat, tt.inputChat)

			// Создаём объект сервиса в который передадим наш мок авторизации
			services := &service.Service{Chat: mockChat}

			// Создаём экземпляр обработчика
			handler := NewHandler(services)

			// Мокируем логгер
			mockLog := slog.New(slog.NewJSONHandler(io.Discard, nil))

			// Инициализируем сервер
			r := chi.NewRouter()
			r.Post("/chats/delete", handler.ChatDelete(mockLog))

			// Готовим тестовый запрос
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/chats/delete", strings.NewReader(tt.inputBody))

			// Выполняем запрос добавляя в контекст userID
			r.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), userCtx, 1)))

			// Сравниваем ожидаемый и актуальный результат
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}

}

func TestHandler_ChatGet(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehavior func(s *mockService.MockChat, chat dto.ChatGet)
	var userID int64 = 1

	// Тестовая таблица с данными
	testTable := []struct {
		name                 string
		inputBody            string
		inputChat            dto.ChatGet
		mockBehavior         mockBehavior
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK one chat",
			inputBody: `{"user_id": 1}`,
			inputChat: dto.ChatGet{
				UserID: &userID,
			},
			mockBehavior: func(s *mockService.MockChat, user dto.ChatGet) {
				s.EXPECT().GetChat(user).Return([]entity.Chat{
					{
						Id:        1,
						Name:      "chat_1",
						CreatedAt: "2024-09-20T18:26:13.239627Z",
						IsDeleted: false,
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Chats get successfully","chats_list":[{"id":1,"name":"chat_1","created_at":"2024-09-20T18:26:13.239627Z","is_deleted":false}]}`,
		},
		{
			name:      "OK many chats",
			inputBody: `{"user_id": 1}`,
			inputChat: dto.ChatGet{
				UserID: &userID,
			},
			mockBehavior: func(s *mockService.MockChat, user dto.ChatGet) {
				s.EXPECT().GetChat(user).Return([]entity.Chat{
					{
						Id:        1,
						Name:      "chat_1",
						CreatedAt: "2024-09-20T18:26:13.239627Z",
						IsDeleted: false,
					},
					{
						Id:        2,
						Name:      "chat_2",
						CreatedAt: "2024-09-19T18:26:13.239627Z",
						IsDeleted: false,
					},
					{
						Id:        3,
						Name:      "chat_1",
						CreatedAt: "2024-09-18T18:26:13.239627Z",
						IsDeleted: false,
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Chats get successfully","chats_list":[{"id":1,"name":"chat_1","created_at":"2024-09-20T18:26:13.239627Z","is_deleted":false},{"id":2,"name":"chat_2","created_at":"2024-09-19T18:26:13.239627Z","is_deleted":false},{"id":3,"name":"chat_1","created_at":"2024-09-18T18:26:13.239627Z","is_deleted":false}]}`,
		},
		{
			name:      "User has no chats",
			inputBody: `{"user_id": 1}`,
			inputChat: dto.ChatGet{
				UserID: &userID,
			},
			mockBehavior: func(s *mockService.MockChat, user dto.ChatGet) {
				s.EXPECT().GetChat(user).Return([]entity.Chat{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"User has no chats"}`,
		},
		{
			name:                 "Required field user_id is missing",
			inputBody:            `{}`,
			mockBehavior:         func(s *mockService.MockChat, user dto.ChatGet) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field UserID is a required field"}`,
		},
		{
			name:                 "Empty user_id",
			inputBody:            `{"user_id":}`,
			mockBehavior:         func(s *mockService.MockChat, user dto.ChatGet) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:                 "Invalid user ID",
			inputBody:            `{"user_id": 0}`,
			mockBehavior:         func(s *mockService.MockChat, user dto.ChatGet) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid user ID"}`,
		},
		{
			name:                 "Error decode",
			inputBody:            `{"user_id"}`,
			mockBehavior:         func(s *mockService.MockChat, user dto.ChatGet) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:                 "Request body is nil",
			mockBehavior:         func(s *mockService.MockChat, user dto.ChatGet) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Empty request"}`,
		},
		{
			name:      "Other error",
			inputBody: `{"user_id": 1}`,
			inputChat: dto.ChatGet{
				UserID: &userID,
			},
			mockBehavior: func(s *mockService.MockChat, user dto.ChatGet) {
				s.EXPECT().GetChat(user).Return(nil, errors.New("some error"))
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Failed to get chats: some error"}`,
		},
	}

	// Итерируемся по нашей тестовой таблице
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			// Инициализируем зависимости

			//Инициализируем контролер для мока сервиса
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Создаём моки сервиса чат
			mockChat := mockService.NewMockChat(ctrl)

			// Передаём структуру пользователя
			tt.mockBehavior(mockChat, tt.inputChat)

			// Создаём объект сервиса в который передадим наш мок авторизации
			services := &service.Service{Chat: mockChat}

			// Создаём экземпляр обработчика
			handler := NewHandler(services)

			// Мокируем логгер
			mockLog := slog.New(slog.NewJSONHandler(io.Discard, nil))

			// Инициализируем сервер
			r := chi.NewRouter()
			r.Post("/chats/get", handler.ChatGet(mockLog))

			// Готовим тестовый запрос
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/chats/get", strings.NewReader(tt.inputBody))

			// Выполняем запрос добавляя в контекст userID
			r.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), userCtx, 1)))

			// Сравниваем ожидаемый и актуальный результат
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}

}
