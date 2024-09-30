package validate

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"service-chat/internal/dto"
	respMock "service-chat/internal/handler/mocks"
)

func TestBaseValidate(t *testing.T) {
	tests := []struct {
		name        string
		log         *slog.Logger
		inputBody   string
		body        func(inputBody string) io.ReadCloser
		req         dto.SignUpRequest
		want        *Result
		errValidate bool
	}{
		{
			name:      "Valid request",
			log:       slog.New(slog.NewJSONHandler(io.Discard, nil)),
			inputBody: `{"username": "Andrey","password": "adgui*"}`,
			body: func(inputBody string) io.ReadCloser {
				return httptest.NewRequest(
					http.MethodPost, "/sign-up", bytes.NewBufferString(inputBody)).Body
			},
		},
		{
			name:      "Empty request",
			log:       slog.New(slog.NewJSONHandler(io.Discard, nil)),
			inputBody: ``,
			body: func(inputBody string) io.ReadCloser {
				return httptest.NewRequest(
					http.MethodPost, "/sign-up", bytes.NewBufferString(inputBody)).Body
			},
			want: &Result{
				ErrMsg: errBody,
			},
		},
		{
			name:      "Invalid request",
			log:       slog.New(slog.NewJSONHandler(io.Discard, nil)),
			inputBody: `{"username": "Andrey","password":}`,
			body: func(inputBody string) io.ReadCloser {
				return httptest.NewRequest(
					http.MethodPost, "/sign-up", bytes.NewBufferString(inputBody)).Body
			},
			want: &Result{
				ErrMsg: errDecode,
			},
		},
		{
			name:      "Incorrect request",
			log:       slog.New(slog.NewJSONHandler(io.Discard, nil)),
			inputBody: `{"username": "Andrey","password": "adgui"}`,
			body: func(inputBody string) io.ReadCloser {
				return httptest.NewRequest(
					http.MethodPost, "/sign-up", bytes.NewBufferString(inputBody)).Body
			},
			want: &Result{
				ValidateErr: validator.ValidationErrors{
					&respMock.MockFieldError{TestTag: "min", TestField: "Password", TestParam: "6"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := tt.body(tt.inputBody)
			actual := BaseValidate(tt.log, body, &tt.req)
			if tt.want != nil && tt.want.ValidateErr != nil {
				assert.Equal(t, tt.want.ValidateErr[0].ActualTag(), actual.ValidateErr[0].ActualTag())
				assert.Equal(t, tt.want.ValidateErr[0].Param(), actual.ValidateErr[0].Param())
				assert.Equal(t, tt.want.ValidateErr[0].Field(), actual.ValidateErr[0].Field())
			} else {
				assert.Equal(t, tt.want, actual)
			}
		})
	}
}
