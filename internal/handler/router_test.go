package handler

import (
	"io"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"service-chat/internal/service"
	mockService "service-chat/internal/service/mocks"
)

func TestNewHandler(t *testing.T) {
	// Определяем функцию для создания мока сервиса
	type mockBehaviour func(ctrl *gomock.Controller) *service.Service

	//Инициализируем контролер для мока сервиса
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Тестовая таблица с данными
	testTable := []struct {
		name     string
		mockServ mockBehaviour
	}{
		{
			name: "All services are correct",
			mockServ: func(ctrl *gomock.Controller) *service.Service {
				return &service.Service{
					Authorization: mockService.NewMockAuthorization(ctrl),
					Chat:          mockService.NewMockChat(ctrl),
					Message:       mockService.NewMockMessage(ctrl),
				}
			},
		},
		{
			name: "Authorization is required",
			mockServ: func(ctrl *gomock.Controller) *service.Service {
				return &service.Service{
					Authorization: mockService.NewMockAuthorization(ctrl),
				}
			},
		},
		{
			name: "Chat is required",
			mockServ: func(ctrl *gomock.Controller) *service.Service {
				return &service.Service{
					Chat: mockService.NewMockChat(ctrl),
				}
			},
		},
		{
			name: "Message is required",
			mockServ: func(ctrl *gomock.Controller) *service.Service {
				return &service.Service{
					Message: mockService.NewMockMessage(ctrl),
				}
			},
		},
		{
			name: "Nil",
			mockServ: func(ctrl *gomock.Controller) *service.Service {
				return &service.Service{}
			},
		},
	}

	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {

			// Создаём моки сервиса
			services := tt.mockServ(ctrl)

			// Создаём экземпляр обработчика
			handler := NewHandler(services)

			// Сравниваем ожидаемый и актуальный результат
			assert.Equalf(t, handler.services, services, "Expected services: %v\n, Actual: %v", handler.services, services)
		})
	}
}

func TestHandler_NewRouter(t *testing.T) {
	// Определяем функцию для создания мока сервиса
	type mockBehaviour func(ctrl *gomock.Controller) *service.Service

	//Инициализируем контролер для мока сервиса
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Тестовая таблица с данными
	testTable := []struct {
		name     string
		mockServ mockBehaviour
		log      *slog.Logger
		patterns []string
	}{
		{
			name: "Simple logger",
			mockServ: func(ctrl *gomock.Controller) *service.Service {
				return &service.Service{
					Authorization: mockService.NewMockAuthorization(ctrl),
					Chat:          mockService.NewMockChat(ctrl),
					Message:       mockService.NewMockMessage(ctrl),
				}
			},
			log:      slog.New(slog.NewJSONHandler(io.Discard, nil)),
			patterns: []string{"/auth/*", "/chats/*", "/messages/*", "/swagger/*"},
		},
	}

	// Итерируемся по нашей тестовой таблице
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			// Создаём моки сервиса
			services := tt.mockServ(ctrl)

			// Создаём экземпляр обработчика
			h := &Handler{
				services: services,
			}

			// Смотрим главные маршрутизаторы
			routes := h.NewRouter(tt.log).Routes()

			// Проверяем, что роутер содержит главные паттерны endPoints
			for _, route := range routes {
				assert.True(t, checkPattern(route.Pattern, tt.patterns))
			}
		})
	}
}

func checkPattern(routePattern string, patterns []string) bool {
	result := false
	for i := 0; i < len(patterns); i++ {
		if routePattern == patterns[i] {
			result = true
			break
		}
	}
	return result
}
