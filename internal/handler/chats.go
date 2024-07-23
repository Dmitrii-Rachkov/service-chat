package handler

import (
	"fmt"
	"net/http"
)

// ChatAdd - создать новый чат между пользователями
// @Summary ChatAdd
// @Security ApiKeyAuth
// @Tags Chat
// @Description Create chat
// @ID Create chat
// @Accept json
// @Produce json
// @Param input body entity.Chat true "chat info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /chats/add [post]
func (h *Handler) ChatAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>ChatAdd</h1>")
		if err != nil {
			return
		}
		_ = fprintf
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
