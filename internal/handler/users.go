package handler

import (
	"fmt"
	"net/http"

	"github.com/go-chi/render"
)

// UserAdd - добавить нового пользователя
// @Summary UserAdd
// @Security ApiKeyAuth
// @Tags User
// @Description Add new user
// @ID Add new user
// @Accept json
// @Produce json
// @Param input body entity.User true "user info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /users/add [post]
func (h *Handler) UserAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := r.Context().Value(userCtx)
		render.JSON(w, r, map[string]interface{}{"id": id})

		//fprintf, err := fmt.Fprintf(w, "<h1>UserAdd</h1>")
		//if err != nil {
		//	return
		//}
		//_ = fprintf
	}
}

// UserUpdate - редактирование пользователя
// @Summary UserUpdate
// @Security ApiKeyAuth
// @Tags User
// @Description Update user
// @ID Update user
// @Accept json
// @Produce json
// @Param input body entity.User true "user info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /users/update [put]
func (h *Handler) UserUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>UserUpdate</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}

// UserDelete - удаление пользователя
// @Summary UserDelete
// @Security ApiKeyAuth
// @Tags User
// @Description Delete user
// @ID Delete user
// @Accept json
// @Produce json
// @Param input body entity.User true "user info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} Response
// @Failure 500 {object} Response
// @Failure default {object} Response
// @Router /users/delete [delete]
func (h *Handler) UserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>UserDelete</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}
