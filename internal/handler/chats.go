package handler

import (
	"fmt"
	"net/http"
)

// ChatAdd - создать новый чат между пользователями
func ChatAdd() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>ChatAdd</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}

// ChatDelete - удалить чат
func ChatDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>ChatDelete</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}

// ChatGet - получить список чатов конкретного пользователя
func ChatGet() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fprintf, err := fmt.Fprintf(w, "<h1>ChatGet</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	}
}
