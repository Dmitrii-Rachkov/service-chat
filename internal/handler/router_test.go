package handler

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"service-chat/internal/service"
	mockService "service-chat/internal/service/mocks"
)

func TestNewHandler(t *testing.T) {
	// Определяем функцию для создания мока сервиса
	type mockBehaviour func(ctrl *gomock.Controller) *service.Service

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
			//Инициализируем контролер для мока сервиса
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			// Создаём моки сервиса
			services := tt.mockServ(ctrl)

			// Создаём экземпляр обработчика
			handler := NewHandler(services)

			// Сравниваем ожидаемый и актуальный результат
			assert.Equalf(t, handler.services, services, "Expected services: %v\n, Actual: %v", handler.services, services)
		})
	}
}
