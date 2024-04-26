package handler

import (
	"fmt"
	"net/http"
)

// MessageAdd - отправить сообщение в чат от лица пользователя
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
func MessageGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>MessageGet</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}
