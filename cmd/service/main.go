package main

import (
	"fmt"
	"net/http"
	"service-chat/internal/config"
)

func main() {
	// Получаем конфиг из файла local.yaml
	cfg := config.MustSetEnv()
	fmt.Println(cfg)

	//// connect to the DB (example)
	//db, err := config.ConnectDB(cfg.Database.Host, cfg.Database.Port, cfg.Database.Username,
	//	cfg.Database.Password, cfg.Database.Name, cfg.Database.Connections)
	//_ = db
	//if err != nil {
	//	fmt.Println("fail connect to db")
	//}
	//
	//http.HandleFunc("/db", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Fprintf(w, "<h1>%s</h1>", cfg.Greeting)
	//})
	//
	//http.ListenAndServe(":"+cfg.Server.Port, nil)

	// Тестовый handler
	http.HandleFunc("/docker", func(writer http.ResponseWriter, request *http.Request) {
		fprintf, err := fmt.Fprintf(writer, "<h1>Hello from Docker container!</h1>")
		if err != nil {
			return
		}
		_ = fprintf
	})

	err := http.ListenAndServe(":"+cfg.Server.Port, nil)
	if err != nil {
		return
	}

	// TODO: init logger: slog
	// библиотека slog является стандартной с версии go 1.21 и она самая актуальная

	// TODO: init storage (db): postgres
	// наверное самая популярная реляционная база данных

	// TODO: init router: chi, chi render
	// удобный, минималистичный, активно развивается, совместим с http/net стандартным пакетом
	// chi render это один из пакетов chi для работы с структурами запросов и ответов

	// TODO: handler sign-up и sign-in
	// реализуем handlers для регистрации и авторизации с помощью JWT

	// TODO: handlers for work service
	// реализуем все остальные handlers для работы основной логики сервиса

	// TODO: schema swagger
	// делаем свагер схему

	// TODO: tests
	// делаем unit тесты

	// TODO: run server
	// просто запускаем сервер
}
