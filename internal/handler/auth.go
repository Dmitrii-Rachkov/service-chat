package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"service-chat/internal/entity"
	"service-chat/internal/logger"
	"service-chat/internal/validate"
)

// SignUp - регистрация пользователя
// @Summary SignUp
// @Tags Auth
// @Description User registration
// @ID User registration
// @Accept json
// @Produce json
// @Param input body entity.User true "user info"
// @Success 200 {object} Response
// @Failure 400,404,405 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /auth/sign-up [post]
func (h *Handler) SignUp(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Логируем наш запрос
		const op = "handler.SignUp"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Структура для записи входных данных из JSON от пользователя
		var req entity.User

		// Анализируем запрос от пользователя
		data := validate.BaseValidate(log, r.Body, req)
		if data.ValidateErr != nil {
			render.JSON(w, r, ValidationError(data.ValidateErr))

			return
		} else if data.Err != "" {
			render.JSON(w, r, Error(data.Err))

			return
		}

		// Отправляем валидную структуру на слой сервиса
		id, errCreate := h.services.Authorization.CreateUser(data.User)
		if errCreate != nil && strings.Contains(errCreate.Error(), "unique_violation") {
			log.Error("user already exists", logger.Err(errCreate))
			render.JSON(w, r, Error("User already exists"))

			return
		} else if errCreate != nil {
			log.Error("failed to create user", logger.Err(errCreate))
			render.JSON(w, r, Error("Failed to create user"))

			return
		}

		// Если ошибок нет отправляем успешный ответ
		log.Info("Create user is successful", slog.Int("id", id))
		render.JSON(w, r, OK(fmt.Sprintf("Create user is successful, id: %d", id)))

		return
	}
}

// SignIn - авторизация пользователя
// @Summary SignIn
// @Tags Auth
// @Description User authorization
// @ID User authorization
// @Accept json
// @Produce json
// @Param input body signInInput true "credentials"
// @Success 200 {string} string "token"
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /auth/sign-in [post]
func (h *Handler) SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>SignIn</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}
