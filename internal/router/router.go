package router

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"service-chat/internal/handler"
)

func NewRouter() *chi.Mux {
	// Инициализируем роутер
	r := chi.NewRouter()

	// Добавляем requestID к каждому запросу
	r.Use(middleware.RequestID)

	// Логируем все входящие запросы
	r.Use(middleware.Logger)

	// Если случается какая-то паника внутри одного из handler,
	// то не должно падать всё приложение, мы восстанавливаем его
	r.Use(middleware.Recoverer)

	// Пишем красивые URL при подключении к обработчику
	r.Use(middleware.URLFormat)

	// Тестовый handler
	r.Get("/docker", func(writer http.ResponseWriter, request *http.Request) {
		fprintf, err := fmt.Fprintf(writer, "<h1>Hello from Docker container!</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	})

	// Структура наших handlers

	// Регистрация и авторизация
	r.Route("/auth", func(r chi.Router) {
		r.Post("/sign-up", handler.SignUp()) // POST /auth/sign-up
		r.Post("/sign-in", handler.SignIn()) // POST /auth/sign-in
	})

	// Работа с сущностью пользователя
	r.Route("/users", func(r chi.Router) {
		r.Post("/add", handler.UserAdd())         // POST /users/add
		r.Put("/update", handler.UserUpdate())    // PUT /users/update
		r.Delete("/delete", handler.UserDelete()) // DELETE /users/delete
	})

	// Работа с чатами
	r.Route("/chats", func(r chi.Router) {
		r.Post("/add", handler.ChatAdd())         // POST /chats/add
		r.Delete("/delete", handler.ChatDelete()) // DELETE /chats/delete
		r.Post("/get", handler.ChatGet())         // POST /chats/get
	})

	// Работа с сообщениями
	r.Route("/messages", func(r chi.Router) {
		r.Post("/add", handler.MessageAdd()) // POST /messages/add
		r.Post("/get", handler.MessageGet()) // POST /messages/get
	})

	return r
}
