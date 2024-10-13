package db

import (
	"errors"
	"strings"
)

const (
	errCodeUnique = "23505"
	errNoRows     = "sql: no rows in result set"
	errChatExist  = "Chat does not exist or has been deleted"
)

func pureErr(err error) error {
	if err == nil {
		return nil
	}

	// Удаляем из текста ошибки обозначение базы данных
	replacer := strings.NewReplacer("pq: ", "", "sql: ", "", "driver: ", "")

	return errors.New(replacer.Replace(err.Error()))
}
