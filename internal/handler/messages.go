package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"service-chat/internal/dto"
	"service-chat/internal/logger"
	"service-chat/internal/validate"
)

// MessageAdd - отправить сообщение в чат от лица пользователя
// @Summary MessageAdd
// @Security ApiKeyAuth
// @Tags Message
// @Description Send message
// @ID Send message
// @Accept json
// @Produce json
// @Param input body dto.MessageAdd true "message info"
// @Success 200 {object} Response
// @Failure 400,404,405 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /messages/add [post]
func (h *Handler) MessageAdd(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Логируем наш запрос
		const op = "handler.MessageAdd"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Получаем id пользователя из контекста
		idCtx, errCtx := GetUserID(r.Context())
		if errCtx != nil {
			log.Error("failed to get userID from context")
			render.JSON(w, r, Error(errCtx.Error()))
		}

		// Структура для записи входных данных из JSON от пользователя
		var req dto.MessageAdd

		// Анализируем запрос от пользователя
		fail := validate.BaseValidate(log, r.Body, &req)
		if fail != nil && fail.ValidateErr != nil {
			log.Error("invalid request data")
			render.JSON(w, r, ValidationError(fail.ValidateErr))
			return
		} else if fail != nil && fail.ErrMsg != "" {
			log.Error("invalid request data")
			render.JSON(w, r, Error(fail.ErrMsg))
			return
		}

		// Проверяем, что id из контекста совпадает с id из запроса
		if int64(idCtx) != req.UserID {
			log.Error("invalid user ID")
			render.JSON(w, r, Error("Invalid user ID"))
			return
		}

		// Отправляем валидную структуру на слой сервиса
		messageID, errMsg := h.services.Message.AddMessage(req)
		if errMsg != nil {
			log.Error("failed to add message", logger.Err(errMsg))
			render.JSON(w, r, Error(fmt.Sprintf("Failed to create chat: %s", errMsg)))
			return
		}

		// Если ошибок нет отправляем успешный ответ
		log.Info("Message created successfully", slog.Int("messageID", messageID))
		render.JSON(w, r, OK(fmt.Sprintf("Message created successfully, id: %d", messageID)))

	}
}

// MessageGet - получить список сообщений в конкретном чате
// @Summary MessageGet
// @Security ApiKeyAuth
// @Tags Message
// @Description Get message
// @ID Get message
// @Accept json
// @Produce json
// @Param input body entity.Message true "message info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /messages/get [post]
func (h *Handler) MessageGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>MessageGet</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}
