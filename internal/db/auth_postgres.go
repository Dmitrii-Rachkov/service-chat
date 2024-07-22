package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"service-chat/internal/entity"
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
	query := `INSERT INTO "user" (username, password_hash) VALUES ($1, $2) RETURNING id`

	// Запрос в базу на создание пользователя
	row := r.db.QueryRow(query, user.Username, user.Password)

	// Если есть ошибка уникальности username
	var rowErr *pq.Error
	ok := errors.As(row.Err(), &rowErr)

	if ok && rowErr.Code == "23505" {
		return 0, fmt.Errorf("error path: %s, error: %s", op, rowErr.Code.Name())
	}

	// Получаем id записи
	if err := row.Scan(&id); err != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", op, err)
	}

	return id, nil
}
