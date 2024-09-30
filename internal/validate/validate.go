package validate

import (
	"errors"
	"io"
	"log/slog"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"service-chat/internal/logger"
)

const (
	errBody     = "Empty request"
	errDecode   = "Invalid request"
	errValidate = "Failed to validate request fields"
)

type Result struct {
	ErrMsg      string
	ValidateErr validator.ValidationErrors
}

func BaseValidate(log *slog.Logger, body io.ReadCloser, req interface{}) *Result {
	// Анализируем запрос от пользователя
	var err error
	err = render.DecodeJSON(body, req)

	// Если получили запрос с пустым телом
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty")

		return &Result{ErrMsg: errBody}
	}

	// Если не удалось декодировать запрос от пользователя
	if err != nil {
		log.Error("failed to decode request body", logger.Err(err))

		return &Result{ErrMsg: errDecode}
	}

	// Проверяем поля запроса на соответствие заданным правилам в dto
	if err = validator.New().Struct(req); err != nil {
		var validateErr validator.ValidationErrors
		match := errors.As(err, &validateErr)
		if !match {
			log.Error("failed to validate required fields in request body", logger.Err(err))

			return &Result{ErrMsg: errValidate}
		}

		log.Error("invalid request", logger.Err(err))

		return &Result{ValidateErr: validateErr}
	}

	// Валидация прошла успешно
	log.Info("request body decoded", slog.Any("requestBody", req))

	return nil
}
