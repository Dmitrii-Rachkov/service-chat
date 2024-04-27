package handler

import (
	"fmt"
	"net/http"
)

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

// SignUp - регистрация пользователя
func SignUp() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>SignUp</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}

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

// SignIn - авторизация пользователя
func SignIn() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>SignIn</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}
