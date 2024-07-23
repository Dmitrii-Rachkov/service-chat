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

// SignUp - регистрация пользователя
// @Summary SignUp
// @Tags Auth
// @Description User registration
// @ID User registration
// @Accept json
// @Produce json
// @Param input body dto.SignUpRequest true "user info"
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
		var req dto.SignUpRequest

		// Анализируем запрос от пользователя
		fail := validate.BaseValidate(log, r.Body, &req)
		if fail != nil && fail.ValidateErr != nil {
			render.JSON(w, r, ValidationError(fail.ValidateErr))

			return
		} else if fail != nil && fail.ErrMsg != "" {
			render.JSON(w, r, Error(fail.ErrMsg))

			return
		}

		// Отправляем валидную структуру на слой сервиса
		id, errCreate := h.services.Authorization.CreateUser(req)
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
// @Param input body dto.SignInRequest true "credentials"
// @Success 200 {object} Response
// @Failure 400,404,405 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /auth/sign-in [post]
func (h *Handler) SignIn(log *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Логируем наш запрос
		const op = "handler.SignIn"
		log = log.With(
			slog.String("op", op),
			slog.String("request_id", middleware.GetReqID(r.Context())),
		)

		// Структура для записи входных данных из JSON от пользователя
		var req dto.SignInRequest

		// Анализируем запрос от пользователя
		fail := validate.BaseValidate(log, r.Body, &req)
		if fail != nil && fail.ValidateErr != nil {
			render.JSON(w, r, ValidationError(fail.ValidateErr))

			return
		} else if fail != nil && fail.ErrMsg != "" {
			render.JSON(w, r, Error(fail.ErrMsg))

			return
		}

		// Отправляем валидную структуру на слой сервиса
		token, errToken := h.services.Authorization.GenerateToken(req)
		if errToken != nil && strings.Contains(errToken.Error(), "sql: no rows in result set") {
			log.Error("user not found", logger.Err(errToken))
			render.JSON(w, r, Error("User not found"))

			return
		} else if errToken != nil {
			log.Error("failed to generation jwt token", logger.Err(errToken))
			render.JSON(w, r, Error(fmt.Sprintf("Failed to generation jwt token: %s", errToken)))

			return
		}

		// Если ошибок нет отправляем успешный ответ
		log.Info("Authorization successful", slog.String("token", token))
		render.JSON(w, r, OK(fmt.Sprintf("Authorization successful, token: %s", token)))

		return
	}
}
