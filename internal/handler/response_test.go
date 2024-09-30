package handler

import (
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	respMock "service-chat/internal/handler/mocks"
)

func TestValidationError(t *testing.T) {
	// Тестовая таблица с данными
	testTable := []struct {
		name       string
		mockResp   validator.ValidationErrors
		exResponse Response
	}{
		{
			name: "Required one field",
			// Используем ручные моки
			mockResp: validator.ValidationErrors{
				&respMock.MockFieldError{TestTag: "required", TestField: "Username"},
			},
			exResponse: Response{
				Status: StatusError,
				Error:  "Field Username is a required field",
			},
		},
		{
			name: "Required two field",
			mockResp: validator.ValidationErrors{
				&respMock.MockFieldError{TestTag: "required", TestField: "Username"},
				&respMock.MockFieldError{TestTag: "required", TestField: "Password"},
			},
			exResponse: Response{
				Status: StatusError,
				Error:  "Field Username is a required field, Field Password is a required field",
			},
		},
		{
			name: "Max elements",
			mockResp: validator.ValidationErrors{
				&respMock.MockFieldError{TestTag: "max", TestField: "Username", TestParam: "20"},
			},
			exResponse: Response{
				Status: StatusError,
				Error:  "Field Username cannot exceed 20 characters",
			},
		},
		{
			name: "Min elements",
			mockResp: validator.ValidationErrors{
				&respMock.MockFieldError{TestTag: "min", TestField: "Password", TestParam: "6"},
			},
			exResponse: Response{
				Status: StatusError,
				Error:  "Field Password must contain at least 6 characters",
			},
		},
		{
			name: "Max and Min elements",
			mockResp: validator.ValidationErrors{
				&respMock.MockFieldError{TestTag: "max", TestField: "Password", TestParam: "12"},
				&respMock.MockFieldError{TestTag: "min", TestField: "Password", TestParam: "6"},
			},
			exResponse: Response{
				Status: StatusError,
				Error:  "Field Password cannot exceed 12 characters, Field Password must contain at least 6 characters",
			},
		},
		{
			name: "Excludes all",
			mockResp: validator.ValidationErrors{
				&respMock.MockFieldError{TestTag: "excludesall", TestField: "chat_name"},
			},
			exResponse: Response{
				Status: StatusError,
				Error:  "Field chat_name must not contain symbols !@#$&*()?",
			},
		},
		{
			name: "Contains any",
			mockResp: validator.ValidationErrors{
				&respMock.MockFieldError{TestTag: "containsany", TestField: "Password"},
			},
			exResponse: Response{
				Status: StatusError,
				Error:  "Field Password must contain Latin letters and Arabic numerals, as well as the symbols @#$&*()",
			},
		},
		{
			name: "Default",
			mockResp: validator.ValidationErrors{
				&respMock.MockFieldError{TestTag: "default", TestField: "Chat"},
			},
			exResponse: Response{
				Status: StatusError,
				Error:  "Field Chat is not valid",
			},
		},
		{
			name: "All errors",
			mockResp: validator.ValidationErrors{
				&respMock.MockFieldError{TestTag: "required", TestField: "Username"},
				&respMock.MockFieldError{TestTag: "max", TestField: "Username", TestParam: "20"},
				&respMock.MockFieldError{TestTag: "min", TestField: "Password", TestParam: "6"},
				&respMock.MockFieldError{TestTag: "excludesall", TestField: "chat_name"},
				&respMock.MockFieldError{TestTag: "containsany", TestField: "Password"},
				&respMock.MockFieldError{TestTag: "default", TestField: "Chat"},
			},
			exResponse: Response{
				Status: StatusError,
				Error:  "Field Username is a required field, Field Username cannot exceed 20 characters, Field Password must contain at least 6 characters, Field chat_name must not contain symbols !@#$&*()?, Field Password must contain Latin letters and Arabic numerals, as well as the symbols @#$&*(), Field Chat is not valid",
			},
		},
	}

	// Итерируемся по нашей тестовой таблице
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			// Отправляем тестовые ошибки в функцию
			acResponse := ValidationError(tt.mockResp)

			// Сравниваем ожидаемый и актуальный результат
			assert.Equal(t, tt.exResponse, acResponse)
		})
	}
}

func TestOK(t *testing.T) {
	// Тестовая таблица с данными
	testTable := []struct {
		name       string
		msg        string
		exResponse Response
	}{
		{
			name: "OK",
			msg:  "Chat created successfully",
			exResponse: Response{
				Status:  StatusOK,
				Message: "Chat created successfully",
			},
		},
		{
			name: "Empty message",
			msg:  "",
			exResponse: Response{
				Status:  StatusOK,
				Message: "",
			},
		},
		{
			name: "Nil message",
			exResponse: Response{
				Status:  StatusOK,
				Message: "",
			},
		},
	}

	// Итерируемся по нашей тестовой таблице
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			// Отправляем тестовые сообщения в функцию
			acResponse := OK(tt.msg)

			// Сравниваем ожидаемый и актуальный результат
			assert.Equal(t, tt.exResponse, acResponse)
		})
	}
}

func TestError(t *testing.T) {
	// Тестовая таблица с данными
	testTable := []struct {
		name       string
		msg        string
		exResponse Response
	}{
		{
			name: "OK",
			msg:  "Failed to create chat",
			exResponse: Response{
				Status: StatusError,
				Error:  "Failed to create chat",
			},
		},
		{
			name: "Empty error message",
			msg:  "",
			exResponse: Response{
				Status: StatusError,
				Error:  "",
			},
		},
		{
			name: "Nil error message",
			exResponse: Response{
				Status:  StatusError,
				Message: "",
			},
		},
	}

	// Итерируемся по нашей тестовой таблице
	for _, tt := range testTable {
		t.Run(tt.name, func(t *testing.T) {
			// Отправляем тестовые сообщения в функцию
			acResponse := Error(tt.msg)

			// Сравниваем ожидаемый и актуальный результат
			assert.Equal(t, tt.exResponse, acResponse)
		})
	}
}
