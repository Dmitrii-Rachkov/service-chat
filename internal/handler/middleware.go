package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/go-chi/render"
)

const (
	authHeader     = "Authorization"
	errEmptyHeader = "Authorization header is empty"
	errInvalidAuth = "Invalid authorization header"
	errEmptyToken  = "Token is empty"
	userCtx        = "userID"
	bearerToken    = "Bearer"
)

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Получаем значение из header авторизации
		header := r.Header.Get(authHeader)
		// Если header пустой
		if header == "" {
			render.JSON(w, r, Error(errEmptyHeader))
			return
		}

		// Проверяем, что значение header валидно
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[0] != bearerToken {
			render.JSON(w, r, Error(errInvalidAuth))
			return
		}

		// Если длина значения token == 0
		if len(headerParts[1]) == 0 {
			render.JSON(w, r, Error(errEmptyToken))
			return
		}

		// Получаем id пользователя из jwt token
		userID, err := h.services.Authorization.ParseToken(headerParts[1])
		if err != nil {
			render.JSON(w, r, Error(err.Error()))
			return
		}

		// Добавляем в контекст id нашего пользователя для передачи в следующие handlers
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), userCtx, userID)))
	})
}

func GetUserID(ctx context.Context) (int, error) {
	// Достаём из контекста userID
	id, ok := ctx.Value(userCtx).(int)
	if !ok {
		return 0, errors.New("user id not found")
	}

	return id, nil
}
