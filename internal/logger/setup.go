package logger

import (
	"log/slog"
	"os"

	"service-chat/internal/logger/prettylog"
)

// slog - это не конкретно logger а некая обёртка, это библиотека для работы с логгерами
// под капотом есть дефолтные логгеры: текстовый (например для локали) и JSON (для отправки например в kibana)
// можно использовать и другие логгеры и также писать свои

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

// SetupLogger - установка logger в зависимости от окружения (локально - текст. на dev или prod - JSON)
// и также уровень логгирования в зависимости от окружения (уровень debug, уровень info)
func SetupLogger(env string) *slog.Logger {
	// Объявляем наш logger библиотеки slog
	var log *slog.Logger

	// В зависимости от окружения устанавливаем наш logger
	switch env {
	case envLocal:
		// Используем красивый текстовый handler для записи в Stdout
		// уровень логгирования debug
		opts := &slog.HandlerOptions{
			Level:     slog.LevelDebug,
			AddSource: false,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == "nothing" {
					return slog.Attr{}
				}
				return a
			},
		}
		log = slog.New(prettylog.NewHandler(opts))
	case envDev:
		// Используем стандартный JSON handler для записи чтобы в дальнейшем передавать логи
		// уровень логгирования debug
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		// Используем стандартный JSON handler для записи чтобы в дальнейшем передавать логи
		// уровень логгирования info
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	// Возвращаем объект логгера
	return log
}
