package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"service-chat/internal/db/entity"
)

type ChatsPostgres struct {
	db *sql.DB
}

func NewChatsPostgres(db *sql.DB) *ChatsPostgres {
	return &ChatsPostgres{db: db}
}

// CreateChat - создаём чат между пользователями
func (c *ChatsPostgres) CreateChat(in entity.ChatAdd) (int, error) {
	const op = "db.CreateChat"
	var chatID int

	// Запускаем транзакцию
	tx, err := c.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", op, err)

	}

	// Скелет sql запроса в базу данных для создания чата в таблице chat
	stmtChat, errChat := c.db.Prepare(`INSERT INTO "chat" (name) VALUES ($1) RETURNING id`)
	if errChat != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", op, errChat)
	}

	// Запрос на создание чата между пользователями
	rowChat := tx.Stmt(stmtChat).QueryRow(in.ChatName)

	// Если есть ошибка уникальности chatName
	var rowErr *pq.Error
	ok := errors.As(rowChat.Err(), &rowErr)
	if ok && rowErr.Code == errCodeUnique {
		return 0, fmt.Errorf("error path: %s, error: %s", op, rowErr.Code.Name())
	}

	// Получаем id chat
	if err = rowChat.Scan(&chatID); err != nil {
		// Откатываем транзакцию в случае ошибки
		errTx := tx.Rollback()
		if errTx != nil {
			return 0, fmt.Errorf("error path: %s, error: %s", op, errTx)
		}
		return 0, fmt.Errorf("error path: %s, error: %w", op, err)
	}

	// Скелет sql запроса в базу данных для связи чата и пользователей в таблице users_chat
	stmtUsersChat, errUsersChat := c.db.Prepare(`INSERT INTO "users_chat" (user_id, chat_id) VALUES ($1, $2)`)
	if errUsersChat != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", op, errUsersChat)
	}

	// Запрос на связь чата с пользователями
	for _, user := range in.Users {
		_, errExec := tx.Stmt(stmtUsersChat).Exec(user, chatID)
		if errExec != nil {
			// Откатываем транзакцию в случае ошибки
			errTx := tx.Rollback()
			if errTx != nil {
				return 0, fmt.Errorf("error path: %s, error: %s", op, errTx)
			}
			return 0, fmt.Errorf("error path: %s, error: %w", op, errExec)
		}
	}

	// Возвращаем id chat
	return chatID, tx.Commit()
}
