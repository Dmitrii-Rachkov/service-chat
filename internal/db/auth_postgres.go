package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"service-chat/internal/db/entity"
)

const (
	errCodeUnique = "23505"
)

type AuthPostgres struct {
	db *sql.DB
}

// NewAuthPostgres - конструктор для работы со слоем базы данных
func NewAuthPostgres(db *sql.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

// CreateUser - реализуем интерфейс функции CreateUser, здесь мы непосредственно записываем пользователя в базу данных
func (r *AuthPostgres) CreateUser(user entity.User) (int, error) {
	const op = "db.CreateUser"
	var id int
	// Скелет sql запроса в базу данных
	stmt, err := r.db.Prepare(`INSERT INTO "user" (username, password_hash) VALUES ($1, $2) RETURNING id`)
	if err != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", op, err)
	}

	// Запрос в базу на создание пользователя
	row := stmt.QueryRow(user.Username, user.Password)

	// Если есть ошибка уникальности username
	var rowErr *pq.Error
	ok := errors.As(row.Err(), &rowErr)
	if ok && rowErr.Code == errCodeUnique {
		return 0, fmt.Errorf("error path: %s, error: %s", op, rowErr.Code.Name())
	}

	// Получаем id записи
	if err = row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", op, err)
	}

	return id, nil
}
