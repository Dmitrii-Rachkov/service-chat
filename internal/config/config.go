package config

import (
	"errors"
	"os"
	"path"
	"runtime"
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
	Env      string   `yaml:"env" env-default:"local" env-description:"Environment"`
	Database Database `yaml:"database"`
	Server   Server   `yaml:"server"`
}

// Database - структура конфига базы данных
type Database struct {
	Host        string `yaml:"host"`
	Port        string `yaml:"port"`
	Username    string `yaml:"username"`
	Password    string `yaml:"password"`
	Name        string `yaml:"name"`
	Connections int    `yaml:"connections"`
	SSLMode     string `yaml:"sslmode"`
}

// Server - структура конфига сервера
type Server struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        string        `yaml:"port"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idleTimeout" env-default:"60s"`
}

// MustSetEnv - функция, которая прочитает файл с конфигом и создаст и заполнит объект Config
func MustSetEnv(configPath string) (*Config, error) {
	// Проверяем существует ли файл с конфигом по указанному пути
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, errors.New("config file doesn't exist")
	}

	// Объект конфига
	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		return nil, errors.New("failed to read config: " + err.Error())
	}

	// Ищем откуда запускается функция
	_, filename, _, _ := runtime.Caller(0)

	// Находим корневую директорию
	dir := path.Join(path.Dir(filename), "../../")

	// Получаем переменные окружения из файла .env
	if err := godotenv.Load(dir + "/.env"); err != nil {
		return nil, errors.New("failed to load .env file: " + err.Error())
	}

	// Достаём пароль из .env
	passFromEnv := os.Getenv("DB_PASSWORD")

	// Расшифровываем и устанавливаем пароль для базы данных
	password, errDec := encryption.Decrypt(passFromEnv)
	if errDec != nil {
		return nil, errors.New("failed to decrypt password: " + errDec.Error())
	}
	cfg.Database.Password = password

	return &cfg, nil
}
