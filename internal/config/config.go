package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config - структура конфига из файла ./config/local.yaml
type Config struct {
	// Тег yaml:"env" определяет какое имя будет у параметра Env в yaml файле если мы оттуда будем считывать данные
	// env-default:"local" - окружение по умолчанию
	// env-description:"Environment" - мы можем получать при необходимости описание параметра
	// env:"DB_USER" - будем получать из переменных окружения, т.к. задали его в compose.yaml
	// env:"DB_PASSWORD" - пароль задан в файле password.txt и путь до него прописан в compose.yaml
	// env:"DB_NAME" - также прописан в compose.yaml
	// yaml:"connections" - кол-во одновременных подключений к базе данных задано в local.yaml
	Env      string `yaml:"env" env-default:"local" env-description:"Environment"`
	Database struct {
		Host        string `yaml:"host" env-description:"Database host"`
		Port        string `yaml:"port" env-description:"Database port"`
		Username    string `env:"DB_USER" env-description:"Database user name"`
		Password    string `env:"DB_PASSWORD" env-description:"Database user password"`
		Name        string `env:"DB_NAME" env-description:"Database name"`
		Connections int    `yaml:"connections" env-description:"Total number of database connections"`
	} `yaml:"database"`
	Server struct {
		Host        string        `yaml:"host" env-description:"Server host" env-default:"localhost"`
		Port        string        `yaml:"port"`
		Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
		IdleTimeout time.Duration `yaml:"idleTimeout" env-default:"60s"`
	} `yaml:"server"`
	Greeting string `env:"GREETING" env-description:"Greeting phrase" env-default:"Hello!"`
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

	return &cfg
}
