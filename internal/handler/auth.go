package handler

import (
	"fmt"
	"net/http"
)

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
