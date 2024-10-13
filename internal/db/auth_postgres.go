package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"service-chat/internal/db/entity"
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

	// Если пустой запрос
	if user.Username == "" || user.Password == "" {
		return 0, fmt.Errorf("error path: %s, error: empty username or password", op)
	}

	var id int
	// Скелет sql запроса в базу данных
	stmt, err := r.db.Prepare(`INSERT INTO "user" (username, password_hash) VALUES ($1, $2) RETURNING id`)
	if err != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", op, pureErr(err))
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
		return 0, fmt.Errorf("error path: %s, error: %w", op, pureErr(err))
	}

	return id, nil
}

func (r *AuthPostgres) GetUser(user entity.User) (*entity.User, error) {
	const op = "db.GetUser"

	// Если пустой username
	if user.Username == "" {
		return nil, fmt.Errorf("error path: %s, error: empty username", op)
	}

	var userDB entity.User
	// Скелет sql запроса в базу данных
	stmt, err := r.db.Prepare(`SELECT id, username, password_hash FROM "user" WHERE username = $1`)
	if err != nil {
		return nil, fmt.Errorf("error path: %s, error: %w", op, pureErr(err))
	}

	// Запрос в базу на получение пользователя
	row := stmt.QueryRow(user.Username)

	// Получаем id, username, password_hash из базы данных
	if err = row.Scan(&userDB.Id, &userDB.Username, &userDB.Password); err != nil {
		return nil, fmt.Errorf("error path: %s, error: %w", op, pureErr(err))
	}

	return &userDB, nil
}
