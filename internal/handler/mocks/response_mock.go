package mocks

import "github.com/go-playground/validator/v10"

// MockFieldError - мок для валидатора
type MockFieldError struct {
	validator.FieldError
	TestTag   string
	TestField string
	TestParam string
}

// ActualTag - мок для тега валидации
func (m *MockFieldError) ActualTag() string {
	return m.TestTag
}

// Field - мок для поля валидации
func (m *MockFieldError) Field() string {
	return m.TestField
}

// Param - мок для параметра валидации
func (m *MockFieldError) Param() string {
	return m.TestParam
}
