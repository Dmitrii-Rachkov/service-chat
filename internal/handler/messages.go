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
			render.JSON(w, r, Error(fmt.Sprintf("Failed to create message: %s", errMsg)))
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
// @Param input body dto.MessageGet true "message info"
// @Success 200 {object} Response{Status, Message, MessagesList}
// @Failure 400,404,405 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /messages/get [post]
func (h *Handler) MessageGet(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Логируем наш запрос
		const op = "handler.MessageGet"
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
		var req dto.MessageGet

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

		// Отправляем валидную структуру на слой сервиса
		messages, errMsg := h.services.Message.GetMessage(req, idCtx)
		if errMsg != nil {
			log.Error("failed to get messages", logger.Err(errMsg))
			render.JSON(w, r, Error(fmt.Sprintf("Failed to get messages: %s", errMsg)))
			return
		}

		// Если в чате нет сообщений
		if len(messages) == 0 {
			log.Info("user don't have messages")
			render.JSON(w, r, OK(fmt.Sprintf("User has no messages in chat with id: %d", req.ChatID)))
		} else {
			// Если ошибок нет и есть сообщения отправляем успешный ответ
			log.Info("Message get successfully", "Messages", messages)
			render.JSON(w, r, Response{
				Status:       StatusOK,
				Message:      "Message get successfully",
				MessagesList: messages},
			)
		}
	}
}

// MessageUpdate - отредактировать сообщение от лица пользователя
// @Summary MessageUpdate
// @Security ApiKeyAuth
// @Tags Message
// @Description Update message
// @ID Update message
// @Accept json
// @Produce json
// @Param input body dto.MessageUpdate true "message info"
// @Success 200 {object} Response
// @Failure 400,404,405 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /messages/update [put]
func (h *Handler) MessageUpdate(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Логируем наш запрос
		const op = "handler.MessageUpdate"
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
		var req dto.MessageUpdate

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
		messageID, errMsg := h.services.Message.UpdateMessage(req)
		if errMsg != nil {
			log.Error("failed to update message", logger.Err(errMsg))
			render.JSON(w, r, Error(fmt.Sprintf("Failed to update message: %s", errMsg)))
			return
		}

		// Если ошибок нет отправляем успешный ответ
		log.Info("Message updated successfully", slog.Int("messageID", messageID))
		render.JSON(w, r, OK(fmt.Sprintf("Message update successfully, id: %d", messageID)))
	}
}
