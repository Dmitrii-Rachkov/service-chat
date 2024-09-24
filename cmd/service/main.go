package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

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
	// Получаем путь до конфига
	configPath := filepath.Join("./", "config", "local.yaml")
	fmt.Println("configPath:", configPath)

	// Получаем конфиг из файла local.yaml
	cfg, errCfg := config.MustSetEnv(configPath)
	if errCfg != nil {
		log.Fatalf("Failed to set env vars: %v", errCfg)
	}

	// Создаём logger
	customLog := logger.SetupLogger(cfg.Env)

	// Накатываем миграцию базы данных
	errMigrateUp := db.MigrateUp(cfg, customLog)
	if errMigrateUp != nil {
		customLog.Error("Failed to up migrate", logger.Err(errMigrateUp))
		os.Exit(1)
	}
	customLog.Info("Migrate Up is successful")

	// Выводим в консоль информацию о запуске нашего приложения, параметры конфига и режиме работы logger
	customLog.Info("Start service-chat", slog.String("env", cfg.Env))
	customLog.Debug("Debug messages is on")

	// Инициализируем базу данных
	database, errDB := db.NewPostgresDB(cfg)
	if errDB != nil {
		customLog.Error("Failed to start database", logger.Err(errDB))
		os.Exit(1)
	}
	customLog.Info("Database initialization was successful")

	// Собираем наши слои проекта
	repos := db.NewDB(database)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	// Инициализируем экземпляр сервера
	srv := new(server.Server)

	// Плавный выход из приложения, перестаём принимать все входящие запросы,
	// при этом закончим обработку всех текущих запросов и операций в базе данных.
	// Для этого запускаем сервер в горутине
	go func() {
		// Запускаем сервер
		if errRunServer := srv.Run(cfg, handlers.NewRouter(customLog)); errRunServer != nil {
			customLog.Error("Failed to start server, error:", logger.Err(errRunServer))
			os.Exit(1)
		}
		customLog.Info("Server started")
	}()

	// Поскольку сервер запускается в горутине, это не блокирует выполнение функции main и поэтому
	// сразу выходим из приложения. Чтобы этого избежать добавим блокировку функции main с помощью
	// канала типа os.Signal
	quit := make(chan os.Signal, 1)

	// Запись в канал будет происходить, когда процесс в котором выполняется наше приложение получит
	// сигнал от системы типа SIGTERM или SIGINT
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	// Строка для чтения из канала, которая будет блокировать выполнение главной горутины main
	<-quit

	// Информируем о завершении работы приложения и выходим из него
	customLog.Info("Server is shutting down")
	if err := srv.ShutDown(context.Background()); err != nil {
		customLog.Error("Failed to shutdown server", logger.Err(err))
	}

	// Закрываем все соединения с бд
	if err := database.Close(); err != nil {
		customLog.Error("Failed to close database", logger.Err(err))
	}
}
