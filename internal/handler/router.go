package handler

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	_ "service-chat/docs"
	"service-chat/internal/service"
)

type Handler struct {
	services *service.Service
}

// NewHandler - Слой обработки запросов
func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) NewRouter(log *slog.Logger) *chi.Mux {
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

	// Swagger handler
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://localhost:9000/swagger/doc.json"), //The url pointing to API definition
	))

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
		r.Post("/sign-up", h.SignUp(log)) // POST /auth/sign-up
		r.Post("/sign-in", h.SignIn(log)) // POST /auth/sign-in
	})

	// Protected Endpoints
	r.Group(func(r chi.Router) {
		r.Use(h.AuthMiddleware)
		// Работа с чатами
		r.Route("/chats", func(r chi.Router) {
			r.Post("/add", h.ChatAdd(log))      // POST /chats/add
			r.Delete("/delete", h.ChatDelete()) // DELETE /chats/delete
			r.Post("/get", h.ChatGet(log))      // POST /chats/get
		})

		// Работа с сообщениями
		r.Route("/messages", func(r chi.Router) {
			r.Post("/add", h.MessageAdd(log))         // POST /messages/add
			r.Post("/get", h.MessageGet(log))         // POST /messages/get
			r.Put("/update", h.MessageUpdate(log))    // PUT /messages/update
			r.Delete("/delete", h.MessageDelete(log)) // DELETE /messages/delete
		})
	})

	return r
}
