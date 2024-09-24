package config

import (
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"
)

func TestMustSetEnv_OK(t *testing.T) {
	// Создаём тестовый конфиг
	cfgOK := &Config{
		Env: "local",
		Database: Database{
			Host:        "clair_postgres",
			Port:        "5432",
			Username:    "postgres",
			Password:    "qwerty",
			Name:        "db_chat",
			Connections: 10,
			SSLMode:     "disable",
		},
		Server: Server{
			Host:        "localhost",
			Port:        "9000",
			Timeout:     time.Second * 5,
			IdleTimeout: time.Second * 60,
		},
	}

	// Создаём тестовый yaml с данными конфига
	yamlCfg, err := yaml.Marshal(&cfgOK)
	if err != nil {
		t.Fatal(err)
	}

	// Ищем откуда запускается тест
	_, filename, _, _ := runtime.Caller(0)

	// Находим директорию /config
	dir := path.Join(path.Dir(filename), ".")

	// Создаём файл для готового конфига
	f, err := os.Create(dir + "/cfg_ok.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Записываем yaml с данными в файл
	_, err = io.WriteString(f, string(yamlCfg))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	// Путь к созданному тестовому конфигу
	cfgOKPath := filepath.Join(dir, "cfg_ok.yaml")

	// Проверяем, что считался конфиг из файла и он верный
	actualCfgOK, err := MustSetEnv(cfgOKPath)
	assert.Nil(t, err)
	assert.Equal(t, cfgOK, actualCfgOK)

	// Удаляем тестовый конфиг
	err = os.Remove(dir + "/cfg_ok.yaml")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMustSetEnv_NoExist(t *testing.T) {
	// Проверяем, что нет указанного конфига
	actualCfgErr, err := MustSetEnv("")
	assert.Nil(t, actualCfgErr)
	assert.Equal(t, errors.New("config file doesn't exist"), err)
}

func TestMustSetEnv_FailedRead(t *testing.T) {
	// Ищем откуда запускается тест
	_, filename, _, _ := runtime.Caller(0)

	// Находим директорию /config
	dir := path.Join(path.Dir(filename), ".")

	// Создаём файл для готового конфига
	_, err := os.Create(dir + "/cfg_fail_read.go")
	if err != nil {
		t.Fatal(err)
	}

	// Путь к созданному тестовому конфигу
	cfgReadPath := filepath.Join(dir, "cfg_fail_read.go")

	// Проверяем, что нельзя считать конфиг
	actualCfgRead, err := MustSetEnv(cfgReadPath)
	assert.Nil(t, actualCfgRead)
	assert.Equal(t, errors.New("failed to read config: file format '.go' doesn't supported by the parser"), err)

	// Удаляем тестовый конфиг
	err = os.Remove(dir + "/cfg_fail_read.go")
	if err != nil {
		t.Fatal(err)
	}
}

func TestMustSetEnv_FailedDecrypt(t *testing.T) {
	// Ищем откуда запускается тест
	_, filename, _, _ := runtime.Caller(0)

	// Находим директорию /config
	dir := path.Join(path.Dir(filename), ".")

	// Устанавливаем временную неверную переменную окружения
	t.Setenv("MY_SECRET", "sdsfsf")

	// Создаём пустой конфиг
	cfg := &Config{}

	// Создаём тестовый yaml с данными конфига
	yamlCfg, err := yaml.Marshal(&cfg)
	if err != nil {
		t.Fatal(err)
	}

	// Создаём файл для готового конфига
	f, err := os.Create(dir + "/cfg_decrypt.yaml")
	if err != nil {
		t.Fatal(err)
	}

	// Записываем yaml с данными в файл
	_, err = io.WriteString(f, string(yamlCfg))
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	// Путь к созданному тестовому конфигу
	cfgDecryptPath := filepath.Join(dir, "cfg_decrypt.yaml")

	// Проверяем, что нельзя считать конфиг
	actualCfg3, err := MustSetEnv(cfgDecryptPath)
	assert.Nil(t, actualCfg3)
	assert.Equal(t, errors.New("failed to decrypt password: encryption.Decrypt: crypto/aes: invalid key size 6"), err)

	err = os.Remove(dir + "/cfg_decrypt.yaml")
	if err != nil {
		t.Fatal(err)
	}
}
