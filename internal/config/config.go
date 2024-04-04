package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"

	"service-chat/internal/encryption"
)

// Config - структура конфига из файла ./config/local.yaml
type Config struct {
	// Тег yaml:"env" определяет какое имя будет у параметра Env в yaml файле если мы оттуда будем считывать данные
	// env-default:"local" - окружение по умолчанию
	// yaml:"connections" - кол-во одновременных подключений к базе данных задано в local.yaml
	Env      string `yaml:"env" env-default:"local" env-description:"Environment"`
	Database struct {
		Host        string `yaml:"host"`
		Port        string `yaml:"port"`
		Username    string `yaml:"username"`
		Password    string `yaml:"password"`
		Name        string `yaml:"name"`
		Connections int    `yaml:"connections"`
		SSLMode     string `yaml:"sslmode"`
	} `yaml:"database"`
	Server struct {
		Host        string        `yaml:"host" env-default:"localhost"`
		Port        string        `yaml:"port"`
		Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
		IdleTimeout time.Duration `yaml:"idleTimeout" env-default:"60s"`
	} `yaml:"server"`
}

// MustSetEnv - функция, которая прочитает файл с конфигом и создаст и заполнит объект Config
func MustSetEnv() *Config {
	// Получаем путь до конфига
	configPath := filepath.Join("./", "config", "local.yaml")
	fmt.Println("configPath:", configPath)

	// Проверяем существует ли файл с конфигом по указанному пути
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Fatalf("config file does not exist in path: %s", configPath)
	}

	// Объект конфига
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("cannot read config: %s", err)
	}

	// Получаем переменные окружения из файла .env
	if err := godotenv.Load(); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	// Расшифровываем и устанавливаем пароль для базы данных
	password, errDec := encryption.Decrypt(os.Getenv("DB_PASSWORD"))
	if errDec != nil {
		log.Fatalf("error decrypting your encrypted text: %s", errDec)
	}
	cfg.Database.Password = password

	return &cfg
}
