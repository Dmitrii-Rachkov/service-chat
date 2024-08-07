package main

import (
	"log/slog"
	"os"

	_ "github.com/lib/pq"

	"service-chat/internal/config"
	"service-chat/internal/db"
	"service-chat/internal/handler"
	"service-chat/internal/logger"
	"service-chat/internal/service"
	"service-chat/server"
)

// @title Service Chat
// @version 1.0
// @description Providing an HTTP API for working with user chats and messages

// @host localhost:9000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	// Получаем конфиг из файла local.yaml
	cfg := config.MustSetEnv()

	// Создаём logger
	log := logger.SetupLogger(cfg.Env)

	// Накатываем миграцию базы данных
	errMigrateUp := db.MigrateUp(cfg, log)
	if errMigrateUp != nil {
		log.Error("Failed to up migrate", logger.Err(errMigrateUp))
		os.Exit(1)
	}
	log.Info("Migrate Up is successful")

	// Выводим в консоль информацию о запуске нашего приложения, параметры конфига и режиме работы logger
	log.Info("Start service-chat", slog.String("env", cfg.Env))
	log.Debug("Debug messages is on")

	// Инициализируем базу данных
	database, errDB := db.NewPostgresDB(cfg)
	if errDB != nil {
		log.Error("Failed to start database", logger.Err(errDB))
		os.Exit(1)
	}
	log.Info("Database initialization was successful")

	// Собираем наши слои проекта
	repos := db.NewDB(database)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	// Инициализируем экземпляр сервера
	srv := new(server.Server)

	// Запускаем сервер
	if errRunServer := srv.Run(cfg, handlers.NewRouter(log)); errRunServer != nil {
		log.Error("Failed to start server, error:", logger.Err(errRunServer))
		os.Exit(1)
	}
	log.Info("Server started")

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
	//http.HandleFunc("/docker", func(writer http.ResponseWriter, request *http.Request) {
	//	fprintf, err := fmt.Fprintf(writer, "<h1>Hello from Docker container!</h1>")
	//	if err != nil {
	//		return
	//	}
	//	_ = fprintf
	//})
	//
	//err := http.ListenAndServe(":"+cfg.Server.Port, nil)
	//if err != nil {
	//	return
	//}

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
