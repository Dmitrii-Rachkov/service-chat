package router

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func NewRouter() *chi.Mux {
	// Инициализируем роутер
	router := chi.NewRouter()

	// Добавляем requestID к каждому запросу
	router.Use(middleware.RequestID)

	// Логируем все входящие запросы
	router.Use(middleware.Logger)

	// Если случается какая-то паника внутри одного из handler,
	// то не должно падать всё приложение, мы восстанавливаем его
	router.Use(middleware.Recoverer)

	// Пишем красивые URL при подключении к обработчику
	router.Use(middleware.URLFormat)

	return router
}
