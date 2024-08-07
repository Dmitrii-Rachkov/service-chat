package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"service-chat/internal/dto"
	"service-chat/internal/logger"
	"service-chat/internal/validate"
)

// ChatAdd - создать новый чат между пользователями
// @Summary ChatAdd
// @Security ApiKeyAuth
// @Tags Chat
// @Description Create chat
// @ID Create chat
// @Accept json
// @Produce json
// @Param input body dto.ChatAdd true "chat info"
// @Success 200 {object} Response
// @Failure 400,404,405 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /chats/add [post]
func (h *Handler) ChatAdd(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Логируем наш запрос
		const op = "handler.ChatAdd"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Структура для записи входных данных из JSON от пользователя
		var req dto.ChatAdd

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
		chatID, err := h.services.Chat.CreateChat(req)
		if err != nil && strings.Contains(err.Error(), "unique_violation") {
			log.Error("chat already exists", logger.Err(err))
			render.JSON(w, r, Error("Chat already exists"))
			return
		} else if err != nil {
			log.Error("failed to create chat", logger.Err(err))
			render.JSON(w, r, Error(fmt.Sprintf("Failed to create chat: %s", err)))
			return
		}

		// Если ошибок нет отправляем успешный ответ
		log.Info("Chat created successfully", slog.Int("chatID", chatID))
		render.JSON(w, r, OK(fmt.Sprintf("Chat created successfully, id: %d", chatID)))
	}
}

// ChatDelete - удалить чат
// @Summary ChatDelete
// @Security ApiKeyAuth
// @Tags Chat
// @Description Delete chat
// @ID Delete chat
// @Accept json
// @Produce json
// @Param input body entity.Chat true "chat info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /chats/delete [delete]
func (h *Handler) ChatDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>ChatDelete</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}

// ChatGet - получить список чатов конкретного пользователя
// @Summary ChatGet
// @Security ApiKeyAuth
// @Tags Chat
// @Description Get chat
// @ID Get chat
// @Accept json
// @Produce json
// @Param input body entity.Chat true "chat info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /chats/get [post]
func (h *Handler) ChatGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>ChatGet</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}
