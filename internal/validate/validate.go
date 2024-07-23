package validate

import (
	"errors"
	"io"
	"log/slog"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"service-chat/internal/entity"
	"service-chat/internal/logger"
)

const (
	errBody     = "Empty request"
	errDecode   = "Invalid request"
	errValidate = "Failed to validate request fields"
	errType     = "Unknown event type"
)

type Result struct {
	Err         string
	ValidateErr validator.ValidationErrors
	User        entity.User
	Chat        entity.Chat
	Message     entity.Message
}

func BaseValidate(log *slog.Logger, body io.ReadCloser, req interface{}) Result {
	// Приводим запрос к нужному типу и анализируем его
	var request Result
	var err error

	switch in := req.(type) {
	case entity.User:
		request.User = in
		err = render.DecodeJSON(body, &request.User)
	case entity.Chat:
		request.Chat = in
		err = render.DecodeJSON(body, &request.Chat)
	case entity.Message:
		request.Message = in
		err = render.DecodeJSON(body, &request.Message)
	default:
		request.Err = errType
		return request
	}

	// Если получили запрос с пустым телом
	if errors.Is(err, io.EOF) {
		log.Error("request body is empty")

		return Result{Err: errBody}
	}

	// Если не удалось декодировать запрос от пользователя
	if err != nil {
		log.Error("failed to decode request body", logger.Err(err))

		return Result{Err: errDecode}
	}

	// Проверяем, что обязательные поля Username и Password заполнены верно
	if err = validator.New().Struct(request); err != nil {
		var validateErr validator.ValidationErrors
		match := errors.As(err, &validateErr)
		if !match {
			log.Error("failed to validate required fields in request body", logger.Err(err))

			return Result{Err: errValidate}
		}

		log.Error("invalid request", logger.Err(err))

		return Result{ValidateErr: validateErr}
	}

	// Валидация прошла успешно
	log.Info("request body decoded", slog.Any("requestBody", req))

	return request
}
