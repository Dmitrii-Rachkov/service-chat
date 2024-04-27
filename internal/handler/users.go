package handler

import (
	"fmt"
	"net/http"
)

// @Summary UserAdd
// @Security ApiKeyAuth
// @Tags User
// @Description Add new user
// @ID Add new user
// @Accept json
// @Produce json
// @Param input body entity.User true "user info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/add [post]

// UserAdd - добавить нового пользователя
func UserAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>UserAdd</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}

// @Summary UserUpdate
// @Security ApiKeyAuth
// @Tags User
// @Description Update user
// @ID Update user
// @Accept json
// @Produce json
// @Param input body entity.User true "user info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/update [put]

// UserUpdate - редактирование пользователя
func UserUpdate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>UserUpdate</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}

// @Summary UserDelete
// @Security ApiKeyAuth
// @Tags User
// @Description Delete user
// @ID Delete user
// @Accept json
// @Produce json
// @Param input body entity.User true "user info"
// @Success 200 {integer} integer 1
// @Failure 400,404 {object} errorResponse
// @Failure 500 {object} errorResponse
// @Failure default {object} errorResponse
// @Router /users/delete [delete]

// UserDelete - удаление пользователя
func UserDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>UserDelete</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}
