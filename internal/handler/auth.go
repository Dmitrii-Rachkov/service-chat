package handler

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"

	"service-chat/internal/entity"
	"service-chat/internal/logger"
)

// SignUp - регистрация пользователя
// @Summary SignUp
// @Tags Auth
// @Description User registration
// @ID User registration
// @Accept json
// @Produce json
// @Param input body entity.User true "user info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
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
		err := render.DecodeJSON(r.Body, &req)

		// Если получили запрос с пустым телом
		if errors.Is(err, io.EOF) {
			log.Error("request body is empty")
			render.JSON(w, r, Error("Empty request"))

			return
		}

		// Если не удалось декодировать запрос от пользователя
		if err != nil {
			log.Error("failed to decode request body", logger.Err(err))
			render.JSON(w, r, Error("Failed to decode request"))

			return
		}

		log.Info("request body decoded", slog.Any("requestBody", req))

		// Проверяем, что обязательные поля Username и Password заполнены верно
		if err = validator.New().Struct(req); err != nil {
			var validateErr validator.ValidationErrors
			match := errors.As(err, &validateErr)
			if !match {
				log.Error("failed to validate required fields in request body", logger.Err(err))
				render.JSON(w, r, Error("Failed to validate request fields"))

				return
			}

			log.Error("invalid request", logger.Err(err))
			render.JSON(w, r, ValidationError(validateErr))

			return
		}

		// Отправляем валидную структуру на слой сервиса
		id, errCreate := h.services.Authorization.CreateUser(req)
		if errCreate != nil && strings.Contains(errCreate.Error(), "unique_violation") {
			log.Error("user already exists", logger.Err(errCreate))
			render.JSON(w, r, Error("User already exists"))

			return
		} else if errCreate != nil {
			log.Error("failed to create user", logger.Err(err))
			render.JSON(w, r, Error("Failed to create user"))

			return
		}

		// Если ошибок нет отправляем успешный ответ
		log.Info("Create user is successful", slog.Int("id", id))
		render.JSON(w, r, OK(fmt.Sprintf("Create user is successful, id: %d", id)))

		return
	}
}

type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
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
