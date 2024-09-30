package handler

import (
	"context"
	"errors"
	"fmt"
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

func TestHandler_MessageAdd(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockService.MockMessage, message dto.MessageAdd)

	// Инициализируем зависимости

	//Инициализируем контролер для мока сервиса
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаём моки сервиса сообщений
	mockMessage := mockService.NewMockMessage(ctrl)

	// Создаём объект сервиса в который передадим наш мок сообщения
	services := &service.Service{Message: mockMessage}

	// Создаём экземпляр обработчика
	handler := NewHandler(services)

	// Мокируем логгер
	mockLog := slog.New(slog.NewJSONHandler(io.Discard, nil))

	// Инициализируем сервер
	r := chi.NewRouter()
	r.Post("/messages/add", handler.MessageAdd(mockLog))

	// Тестовая таблица с данными
	testTable := []struct {
		name                 string
		inputBody            string
		inputMessage         dto.MessageAdd
		mockBehaviour        mockBehaviour
		unauthorized         bool
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"chat_id": 1,"user_id": 1,"text": "msg1"}`,
			inputMessage: dto.MessageAdd{
				ChatID: 1,
				UserID: 1,
				Text:   "msg1",
			},
			mockBehaviour: func(s *mockService.MockMessage, message dto.MessageAdd) {
				s.EXPECT().AddMessage(message).Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Message created successfully, id: 1"}`,
		},
		{
			name:                 "Required field chat_id is missing",
			inputBody:            `{"user_id": 1,"text": "msg1"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field ChatID is a required field"}`,
		},
		{
			name:                 "Required field user_id is missing",
			inputBody:            `{"chat_id": 1,"text": "msg1"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field UserID is a required field"}`,
		},
		{
			name:                 "Required field text is missing",
			inputBody:            `{"chat_id": 1,"user_id": 1}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Text is a required field"}`,
		},
		{
			name:                 "Empty chat_id",
			inputBody:            `{"chat_id": ,"user_id": 1,"text": "msg1"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:                 "Empty user_id",
			inputBody:            `{"chat_id": 1,"user_id": ,"text": "msg1"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:                 "Empty text",
			inputBody:            `{"chat_id": 1,"user_id": 1,"text": }`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:                 "Zero chat_id",
			inputBody:            `{"chat_id": 0,"user_id": 1,"text": "msg1"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field ChatID is a required field"}`,
		},
		{
			name:                 "Zero user_id",
			inputBody:            `{"chat_id": 1,"user_id": 0,"text": "msg1"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field UserID is a required field"}`,
		},
		{
			name:                 "Nil text",
			inputBody:            `{"chat_id": 1,"user_id": 1,"text": ""}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Text is a required field"}`,
		},
		{
			name:                 "Request body is nil",
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Empty request"}`,
		},
		{
			name:                 "Error decode",
			inputBody:            `{"chat_id": 1,"user_id": 1,"text":}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:      "Other error",
			inputBody: `{"chat_id": 1,"user_id": 1,"text": "msg1"}`,
			inputMessage: dto.MessageAdd{
				ChatID: 1,
				UserID: 1,
				Text:   "msg1",
			},
			mockBehaviour: func(s *mockService.MockMessage, message dto.MessageAdd) {
				s.EXPECT().AddMessage(message).Return(0, errors.New("some error"))
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Failed to create message: some error"}`,
		},
		{
			name:                 "Invalid user ID",
			inputBody:            `{"chat_id": 1,"user_id": 2,"text": "msg1"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageAdd) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid user ID"}`,
		},
		{
			name:                 "User id not found",
			inputBody:            `{"chat_id": 1,"user_id": 2,"text": "msg1"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageAdd) {},
			unauthorized:         true,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"user id not found"}`,
		},
	}

	// Итерируемся по нашей тестовой таблице
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			// Передаём структуру сообщения
			tt.mockBehaviour(mockMessage, tt.inputMessage)

			// Готовим тестовый запрос
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/messages/add", strings.NewReader(tt.inputBody))

			// Выполняем запрос
			if tt.unauthorized {
				r.ServeHTTP(w, req)
			} else {
				r.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), userCtx, 1)))
			}

			// Сравниваем ожидаемый и актуальный результат
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestHandler_MessageGet(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockService.MockMessage, message dto.MessageGet)
	var (
		limit  = int64(10)
		offset = int64(0)
	)

	// Инициализируем зависимости

	//Инициализируем контролер для мока сервиса
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаём моки сервиса сообщений
	mockMessage := mockService.NewMockMessage(ctrl)

	// Создаём объект сервиса в который передадим наш мок сообщения
	services := &service.Service{Message: mockMessage}

	// Создаём экземпляр обработчика
	handler := NewHandler(services)

	// Мокируем логгер
	mockLog := slog.New(slog.NewJSONHandler(io.Discard, nil))

	// Инициализируем сервер
	r := chi.NewRouter()
	r.Post("/messages/get", handler.MessageGet(mockLog))

	// Тестовая таблица с данными
	testTable := []struct {
		name                 string
		inputBody            string
		inputMessage         dto.MessageGet
		userID               int
		mockBehaviour        mockBehaviour
		unauthorized         bool
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK one message",
			inputBody: `{"chat_id": 1,"limit": 10,"offset": 0}`,
			inputMessage: dto.MessageGet{
				ChatID: 1,
				Limit:  &limit,
				Offset: &offset,
			},
			userID: 1,
			mockBehaviour: func(s *mockService.MockMessage, message dto.MessageGet) {
				s.EXPECT().GetMessage(message, 1).Return([]entity.Message{
					{
						Id:        1,
						Text:      "msg1",
						UserID:    1,
						CreatedAt: "2024-09-20T18:26:13.239627Z",
						IsDeleted: false,
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Message get successfully","messages_list":[{"id":1,"text":"msg1","user_id":1,"created_at":"2024-09-20T18:26:13.239627Z","is_deleted":false}]}`,
		},
		{
			name:      "OK one message",
			inputBody: `{"chat_id": 1,"limit": 10,"offset": 0}`,
			inputMessage: dto.MessageGet{
				ChatID: 1,
				Limit:  &limit,
				Offset: &offset,
			},
			userID: 1,
			mockBehaviour: func(s *mockService.MockMessage, message dto.MessageGet) {
				s.EXPECT().GetMessage(message, 1).Return([]entity.Message{
					{
						Id:        1,
						Text:      "msg1",
						UserID:    1,
						CreatedAt: "2024-09-20T18:26:13.239627Z",
						IsDeleted: false,
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Message get successfully","messages_list":[{"id":1,"text":"msg1","user_id":1,"created_at":"2024-09-20T18:26:13.239627Z","is_deleted":false}]}`,
		},
		{
			name:      "OK many message",
			inputBody: `{"chat_id": 1,"limit": 10,"offset": 0}`,
			inputMessage: dto.MessageGet{
				ChatID: 1,
				Limit:  &limit,
				Offset: &offset,
			},
			userID: 1,
			mockBehaviour: func(s *mockService.MockMessage, message dto.MessageGet) {
				s.EXPECT().GetMessage(message, 1).Return([]entity.Message{
					{
						Id:        1,
						Text:      "msg1",
						UserID:    1,
						CreatedAt: "2024-09-18T18:26:13.239627Z",
						IsDeleted: false,
					},
					{
						Id:        1,
						Text:      "msg2",
						UserID:    1,
						CreatedAt: "2024-09-19T18:26:13.239627Z",
						IsDeleted: false,
					},
					{
						Id:        1,
						Text:      "msg3",
						UserID:    1,
						CreatedAt: "2024-09-20T18:26:13.239627Z",
						IsDeleted: false,
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Message get successfully","messages_list":[{"id":1,"text":"msg1","user_id":1,"created_at":"2024-09-18T18:26:13.239627Z","is_deleted":false},{"id":1,"text":"msg2","user_id":1,"created_at":"2024-09-19T18:26:13.239627Z","is_deleted":false},{"id":1,"text":"msg3","user_id":1,"created_at":"2024-09-20T18:26:13.239627Z","is_deleted":false}]}`,
		},
		{
			name:      "User don't have messages",
			inputBody: `{"chat_id": 1,"limit": 10,"offset": 0}`,
			inputMessage: dto.MessageGet{
				ChatID: 1,
				Limit:  &limit,
				Offset: &offset,
			},
			userID: 1,
			mockBehaviour: func(s *mockService.MockMessage, message dto.MessageGet) {
				s.EXPECT().GetMessage(message, 1).Return([]entity.Message{}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"User has no messages in chat with id: 1"}`,
		},
		{
			name:                 "Required field chat_id is missing",
			inputBody:            `{"limit": 10,"offset": 0}`,
			userID:               1,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageGet) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field ChatID is a required field"}`,
		},
		{
			name:                 "Required field limit is missing",
			inputBody:            `{"chat_id": 1,"offset": 0}`,
			userID:               1,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageGet) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Limit is a required field"}`,
		},
		{
			name:                 "Required field offset is missing",
			inputBody:            `{"chat_id": 1,"limit": 10}`,
			userID:               1,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageGet) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field Offset is a required field"}`,
		},
		{
			name:                 "Empty chat_id",
			inputBody:            `{"chat_id": ,"limit": 10,"offset": 0}`,
			userID:               1,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageGet) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:                 "Empty limit",
			inputBody:            `{"chat_id": 1,"limit": ,"offset": 0}`,
			userID:               1,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageGet) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:                 "Empty offset",
			inputBody:            `{"chat_id": 1,"limit": 10,"offset": }`,
			userID:               1,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageGet) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:                 "Zero chat_id",
			inputBody:            `{"chat_id": 0,"limit": 10,"offset": 0}`,
			userID:               1,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageGet) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field ChatID is a required field"}`,
		},
		{
			name:                 "Request body is nil",
			userID:               1,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageGet) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Empty request"}`,
		},
		{
			name:      "Other error",
			inputBody: `{"chat_id": 1,"limit": 10,"offset": 0}`,
			inputMessage: dto.MessageGet{
				ChatID: 1,
				Limit:  &limit,
				Offset: &offset,
			},
			userID: 1,
			mockBehaviour: func(s *mockService.MockMessage, message dto.MessageGet) {
				s.EXPECT().GetMessage(message, 1).Return(nil, errors.New("some error"))
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Failed to get messages: some error"}`,
		},
		{
			name:      "User not found",
			inputBody: `{"chat_id": 1,"limit": 10,"offset": 0}`,
			inputMessage: dto.MessageGet{
				ChatID: 1,
				Limit:  &limit,
				Offset: &offset,
			},
			userID:               1,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageGet) {},
			unauthorized:         true,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"user id not found"}`,
		},
	}

	// Итерируемся по нашей тестовой таблице
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			// Передаём структуру сообщения
			tt.mockBehaviour(mockMessage, tt.inputMessage)

			// Готовим тестовый запрос
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/messages/get", strings.NewReader(tt.inputBody))

			// Выполняем запрос
			if tt.unauthorized {
				r.ServeHTTP(w, req)
			} else {
				r.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), userCtx, tt.userID)))
			}

			// Сравниваем ожидаемый и актуальный результат
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestHandler_MessageUpdate(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockService.MockMessage, message dto.MessageUpdate)

	// Инициализируем зависимости

	//Инициализируем контролер для мока сервиса
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаём моки сервиса сообщений
	mockMessage := mockService.NewMockMessage(ctrl)

	// Создаём объект сервиса в который передадим наш мок сообщения
	services := &service.Service{Message: mockMessage}

	// Создаём экземпляр обработчика
	handler := NewHandler(services)

	// Мокируем логгер
	mockLog := slog.New(slog.NewJSONHandler(io.Discard, nil))

	// Инициализируем сервер
	r := chi.NewRouter()
	r.Put("/messages/update", handler.MessageUpdate(mockLog))

	// Тестовая таблица с данными
	testTable := []struct {
		name                 string
		inputBody            string
		inputMessage         dto.MessageUpdate
		mockBehaviour        mockBehaviour
		unauthorized         bool
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK",
			inputBody: `{"message_id": 1,"user_id": 1,"new_text": "new_text"}`,
			inputMessage: dto.MessageUpdate{
				MessageID: 1,
				UserID:    1,
				NewText:   "new_text",
			},
			mockBehaviour: func(s *mockService.MockMessage, message dto.MessageUpdate) {
				s.EXPECT().UpdateMessage(message).Return(1, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Message update successfully, id: 1"}`,
		},
		{
			name:                 "Required field message_id is missing",
			inputBody:            `{"user_id": 1,"new_text": "new_text"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageUpdate) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field MessageID is a required field"}`,
		},
		{
			name:                 "Required field user_id is missing",
			inputBody:            `{"message_id": 1,"new_text": "new_text"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageUpdate) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field UserID is a required field"}`,
		},
		{
			name:                 "Required field new_text is missing",
			inputBody:            `{"message_id": 1,"user_id": 1}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageUpdate) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field NewText is a required field"}`,
		},
		{
			name:                 "Zero message_id",
			inputBody:            `{"message_id": 0,"user_id": 1,"new_text": "new_text"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageUpdate) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field MessageID is a required field"}`,
		},
		{
			name:                 "Zero user_id",
			inputBody:            `{"message_id": 1,"user_id": 0,"new_text": "new_text"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageUpdate) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field UserID is a required field"}`,
		},
		{
			name:                 "Nil new_text",
			inputBody:            `{"message_id": 1,"user_id": 1,"new_text": ""}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageUpdate) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field NewText is a required field"}`,
		},
		{
			name:                 "Empty message_id",
			inputBody:            `{"message_id": ,"user_id": 1,"new_text": "new_text"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageUpdate) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:                 "Empty user_id",
			inputBody:            `{"message_id": 1,"user_id": ,"new_text": "new_text"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageUpdate) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:                 "Nil new_text",
			inputBody:            `{"message_id": 1,"user_id": 1,"new_text": }`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageUpdate) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:                 "Request body is nil",
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageUpdate) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Empty request"}`,
		},
		{
			name:      "Other error",
			inputBody: `{"message_id": 1,"user_id": 1,"new_text": "new_text"}`,
			inputMessage: dto.MessageUpdate{
				MessageID: 1,
				UserID:    1,
				NewText:   "new_text",
			},
			mockBehaviour: func(s *mockService.MockMessage, message dto.MessageUpdate) {
				s.EXPECT().UpdateMessage(message).Return(0, errors.New("some error"))
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Failed to update message: some error"}`,
		},
		{
			name:                 "User id not found",
			inputBody:            `{"message_id": 1,"user_id": 1,"new_text": "new_text"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageUpdate) {},
			unauthorized:         true,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"user id not found"}`,
		},
		{
			name:                 "Invalid user ID",
			inputBody:            `{"message_id": 1,"user_id": 2,"new_text": "new_text"}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageUpdate) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid user ID"}`,
		},
	}

	// Итерируемся по нашей тестовой таблице
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			// Передаём структуру сообщения
			tt.mockBehaviour(mockMessage, tt.inputMessage)

			// Готовим тестовый запрос
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPut, "/messages/update", strings.NewReader(tt.inputBody))

			// Выполняем запрос
			if tt.unauthorized {
				r.ServeHTTP(w, req)
			} else {
				r.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), userCtx, 1)))
			}

			// Сравниваем ожидаемый и актуальный результат
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}

func TestHandler_MessageDelete(t *testing.T) {
	// Структура для последующей реализации поведения мока
	type mockBehaviour func(s *mockService.MockMessage, message dto.MessageDelete)

	// Инициализируем зависимости

	//Инициализируем контролер для мока сервиса
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Создаём моки сервиса сообщений
	mockMessage := mockService.NewMockMessage(ctrl)

	// Создаём объект сервиса в который передадим наш мок сообщения
	services := &service.Service{Message: mockMessage}

	// Создаём экземпляр обработчика
	handler := NewHandler(services)

	// Мокируем логгер
	mockLog := slog.New(slog.NewJSONHandler(io.Discard, nil))

	// Инициализируем сервер
	r := chi.NewRouter()
	r.Delete("/messages/delete", handler.MessageDelete(mockLog))

	// Тестовая таблица с данными
	testTable := []struct {
		name                 string
		inputBody            string
		inputMessage         dto.MessageDelete
		mockBehaviour        mockBehaviour
		unauthorized         bool
		expectedStatusCode   int
		expectedResponseBody string
	}{
		{
			name:      "OK one message",
			inputBody: `{"message_ids": [1]}`,
			inputMessage: dto.MessageDelete{
				MessageIds: &[]int64{1},
			},
			mockBehaviour: func(s *mockService.MockMessage, message dto.MessageDelete) {
				s.EXPECT().DeleteMessage(message, 1).Return([]entity.DelMsg{
					{
						MessageID: 1,
						Result:    "Message successfully deleted",
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Result of deleted messages","del_msg_list":[{"message_id":1,"result":"Message successfully deleted"}]}`,
		},
		{
			name:      "OK many message",
			inputBody: `{"message_ids": [1,2,3]}`,
			inputMessage: dto.MessageDelete{
				MessageIds: &[]int64{1, 2, 3},
			},
			mockBehaviour: func(s *mockService.MockMessage, message dto.MessageDelete) {
				s.EXPECT().DeleteMessage(message, 1).Return([]entity.DelMsg{
					{
						MessageID: 1,
						Result:    "Message successfully deleted",
					},
					{
						MessageID: 2,
						Result:    "Message successfully deleted",
					},
					{
						MessageID: 3,
						Result:    "Message successfully deleted",
					},
				}, nil)
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"OK","message":"Result of deleted messages","del_msg_list":[{"message_id":1,"result":"Message successfully deleted"},{"message_id":2,"result":"Message successfully deleted"},{"message_id":3,"result":"Message successfully deleted"}]}`,
		},
		{
			name:                 "Required field message_ids is missing",
			inputBody:            `{}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageDelete) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field MessageIds is a required field"}`,
		},
		{
			name:                 "Empty value message_ids",
			inputBody:            `{"message_ids": []}`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageDelete) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Field MessageIds must contain at least 1 characters"}`,
		},
		{
			name:      "Zero value message_ids",
			inputBody: `{"message_ids": [0]}`,
			inputMessage: dto.MessageDelete{
				MessageIds: &[]int64{0},
			},
			mockBehaviour: func(s *mockService.MockMessage, message dto.MessageDelete) {
				s.EXPECT().DeleteMessage(message, 1).Return(nil, errors.New(
					fmt.Sprintf("Failed to delete message: error path: db.DeleteMessage, error: pq: Not found messages")))
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Failed to delete message: Failed to delete message: error path: db.DeleteMessage, error: pq: Not found messages"}`,
		},
		{
			name:                 "Error decode",
			inputBody:            `{"message_ids": }`,
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageDelete) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Invalid request"}`,
		},
		{
			name:                 "Request nil body",
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageDelete) {},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Empty request"}`,
		},
		{
			name:      "Other error",
			inputBody: `{"message_ids": [1]}`,
			inputMessage: dto.MessageDelete{
				MessageIds: &[]int64{1},
			},
			mockBehaviour: func(s *mockService.MockMessage, message dto.MessageDelete) {
				s.EXPECT().DeleteMessage(message, 1).Return(nil, errors.New("some error"))
			},
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"Failed to delete message: some error"}`,
		},
		{
			name:      "User id not found",
			inputBody: `{"message_ids": [1]}`,
			inputMessage: dto.MessageDelete{
				MessageIds: &[]int64{1},
			},
			mockBehaviour:        func(s *mockService.MockMessage, message dto.MessageDelete) {},
			unauthorized:         true,
			expectedStatusCode:   http.StatusOK,
			expectedResponseBody: `{"status":"Error","error":"user id not found"}`,
		},
	}

	// Итерируемся по нашей тестовой таблице
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			// Передаём структуру сообщения
			tt.mockBehaviour(mockMessage, tt.inputMessage)

			// Готовим тестовый запрос
			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodDelete, "/messages/delete", strings.NewReader(tt.inputBody))

			// Выполняем запрос
			if tt.unauthorized {
				r.ServeHTTP(w, req)
			} else {
				r.ServeHTTP(w, req.WithContext(context.WithValue(req.Context(), userCtx, 1)))
			}

			// Сравниваем ожидаемый и актуальный результат
			assert.Equal(t, tt.expectedStatusCode, w.Code)
			assert.Equal(t, tt.expectedResponseBody, strings.TrimSpace(w.Body.String()))
		})
	}
}
