package handler

import (
	"fmt"
	"net/http"
)

// MessageAdd - отправить сообщение в чат от лица пользователя
// @Summary MessageAdd
// @Security ApiKeyAuth
// @Tags Message
// @Description Send message
// @ID Send message
// @Accept json
// @Produce json
// @Param input body entity.Message true "message info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /messages/add [post]
func MessageAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>MessageAdd</h1>")
		if err != nil {
			return
		}
		_ = fprintf
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
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /messages/get [post]
func MessageGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>MessageGet</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}
