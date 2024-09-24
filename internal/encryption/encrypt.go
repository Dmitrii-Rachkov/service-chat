package encryption

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/joho/godotenv"
)

// Случайный набор байтов
var bytes = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

// Достаём ключ шифрования из переменных окружения
func mySecret() (string, error) {
	// Ищем откуда запускается функция
	_, filename, _, _ := runtime.Caller(0)

	// Находим директорию корневую директорию
	dir := path.Join(path.Dir(filename), "../../")

	// Получаем переменные окружения из файла .env
	if err := godotenv.Load(dir + "/.env"); err != nil {
		return "", errors.New(
			fmt.Sprintf("error loading env variables: %s", err.Error()))
	}

	return os.Getenv("MY_SECRET"), nil
}

// Кодируем и возвращаем строку в base64
func encode(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

// Encrypt - шифруем текст
func Encrypt(text string) (string, error) {
	const op = "encryption.Encrypt"
	secret, err := mySecret()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	plainText := []byte(text)
	cfb := cipher.NewCFBEncrypter(block, bytes)
	cipherText := make([]byte, len(plainText))
	cfb.XORKeyStream(cipherText, plainText)
	return encode(cipherText), nil
}

// Декодируем из base64
func decode(s string) ([]byte, error) {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, errors.New(
			fmt.Sprintf("error DecodeString: %s", err.Error()))
	}
	return data, nil
}

// Decrypt - расшифровываем текст
func Decrypt(text string) (string, error) {
	const op = "encryption.Decrypt"
	secret, err := mySecret()
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	block, err := aes.NewCipher([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	cipherText, err := decode(text)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}
	cfb := cipher.NewCFBDecrypter(block, bytes)
	plainText := make([]byte, len(cipherText))
	cfb.XORKeyStream(plainText, cipherText)

	return string(plainText), nil
}
