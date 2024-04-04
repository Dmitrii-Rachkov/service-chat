package main

import (
	"fmt"
	"log/slog"
	"net/http"

	_ "github.com/lib/pq"

	"service-chat/internal/config"
	"service-chat/internal/db"
	"service-chat/internal/logger"
)

func main() {
	// Получаем конфиг из файла local.yaml
	cfg := config.MustSetEnv()

	// Создаём logger
	log := logger.SetupLogger(cfg.Env)

	// Выводим в консоль информацию о запуске нашего приложения, параметры конфига и режиме работы logger
	log.Info("Start service-chat", slog.String("env", cfg.Env))
	log.Debug("Debug messages is on")

	// Инициализируем базу данных
	database, errDB := db.NewPostgresDB(cfg)
	if errDB != nil {
		log.Error("Failed to start database, error:", errDB.Error())
	}
	log.Info("Database initialization was successful")

	_ = database

	//repos := db.NewDB(database)
	//services := service.NewService(repos)

	// connect to the DB (example)
	//database, errDB := db.Conn(cfg.Database.Host, cfg.Database.Port, cfg.Database.Username,
	//	cfg.Database.Password, cfg.Database.Name, cfg.Database.Connections)
	//_ = database
	//if errDB != nil {
	//	log.Error("Failed to start database, error:", errDB.Error())
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

	// TODO: init server and start

	//// Инициализируем экземпляр сервера
	//srv := new(server.Server)
	//
	//// Запускаем сервер
	//if err = srv.Run(cfg); err != nil {
	//	log.Error("Failed to start server, error:", err.Error())
	//}

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
