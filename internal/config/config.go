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
	// Тег yaml определяет какое имя будет у параметра Env в yaml файле если мы оттуда будем считывать данные
	// envDefault - окружение по умолчанию
	Env          string `yaml:"env" env-default:"local"`
	StoragePaths string `yaml:"storagePath" env-required:"true"`
	HTTPServer   `yaml:"httpServer"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:9000"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idleTimeout" env-default:"60s"`
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
