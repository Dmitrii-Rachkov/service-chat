package handler

import (
	"fmt"
	"net/http"
)

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
