package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"service-chat/internal/db/entity"
)

const (
	opCreateChat = "db.CreateChat"
	opGetChat    = "db.GetChat"
)

type ChatsPostgres struct {
	db *sql.DB
}

func NewChatsPostgres(db *sql.DB) *ChatsPostgres {
	return &ChatsPostgres{db: db}
}

// CreateChat - создаём чат между пользователями
func (c *ChatsPostgres) CreateChat(in entity.ChatAdd) (int, error) {
	var chatID int

	// Запускаем транзакцию
	tx, err := c.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", opCreateChat, err)

	}

	// Скелет sql запроса в базу данных для создания чата в таблице chat
	stmtChat, errChat := c.db.Prepare(`INSERT INTO "chat" (name) VALUES ($1) RETURNING id`)
	if errChat != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", opCreateChat, errChat)
	}
	defer stmtChat.Close()

	// Запрос на создание чата между пользователями
	rowChat := tx.Stmt(stmtChat).QueryRow(in.ChatName)

	// Если есть ошибка уникальности chatName
	var rowErr *pq.Error
	ok := errors.As(rowChat.Err(), &rowErr)
	if ok && rowErr.Code == errCodeUnique {
		return 0, fmt.Errorf("error path: %s, error: %s", opCreateChat, rowErr.Code.Name())
	}

	// Получаем id chat
	if err = rowChat.Scan(&chatID); err != nil {
		// Откатываем транзакцию в случае ошибки
		errTx := tx.Rollback()
		if errTx != nil {
			return 0, fmt.Errorf("error path: %s, error: %s", opCreateChat, errTx)
		}
		return 0, fmt.Errorf("error path: %s, error: %w", opCreateChat, err)
	}

	// Скелет sql запроса в базу данных для связи чата и пользователей в таблице users_chat
	stmtUsersChat, errUsersChat := tx.Prepare(pq.CopyIn("users_chat", "user_id", "chat_id"))
	if errUsersChat != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", opCreateChat, errUsersChat)
	}
	defer stmtUsersChat.Close()

	// Запрос на связь чата с пользователями
	for _, user := range in.Users {
		_, err = stmtUsersChat.Exec(user, chatID)
		if err != nil {
			// Откатываем транзакцию в случае ошибки
			errTx := tx.Rollback()
			if errTx != nil {
				return 0, fmt.Errorf("error path: %s, error: %s", opCreateChat, errTx)
			}
			// Возвращаем ошибку
			return 0, fmt.Errorf("error path: %s, error: %w", opCreateChat, err)
		}
	}

	// Очищаем все буферизованные данные
	_, err = stmtUsersChat.Exec()
	if err != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", opCreateChat, err)
	}

	// Возвращаем id chat и завершаем транзакцию
	return chatID, tx.Commit()
}

// GetChat - получаем все чаты пользователя из бд
func (c *ChatsPostgres) GetChat(in entity.ChatGet) ([]entity.Chat, error) {
	// Скелет sql запроса на получение всех чатов из бд для конкретного пользователя
	// в зависимости от последнего сообщения в каждом из чатов

	// Если в чатах нет сообщений, то выводим все чаты сортирую их по дате создания чата [Z-A]

	// Если в одних чатах есть сообщения, а в других нет, то сначала выводим чаты с сообщениями, сортируя от [Z-A],
	// затем выводим пустые чаты с сортировкой по дате создания чата от [Z-A]
	stmtChats, err := c.db.Prepare(`WITH sort_chat AS (
												SELECT uc.chat_id, MAX(cm.message_id) AS mm, c.id, c.name, c.created_at, c.is_deleted
												FROM users_chat AS uc
												LEFT OUTER JOIN chats_messages AS cm
												ON uc.id = cm.users_chat_id
												RIGHT JOIN chat AS c
												ON c.id = uc.chat_id
												WHERE uc.user_id = $1
												GROUP BY uc.chat_id, c.id
												ORDER BY mm DESC NULLS LAST, uc.chat_id DESC
											)
											SELECT sort_chat.id AS id, sort_chat.name, sort_chat.created_at, sort_chat.is_deleted
											FROM sort_chat`)

	if err != nil {
		return nil, fmt.Errorf("error path: %s, error: %w", opGetChat, err)
	}
	defer stmtChats.Close()

	// Получаем чаты из бд
	rowsChats, err := stmtChats.Query(in.UserID)
	if err != nil {
		return nil, fmt.Errorf("error path: %s, error: %w", opGetChat, err)
	}
	defer rowsChats.Close()

	// Структура для записи всех полученных чатов из бд
	var chats []entity.Chat
	for rowsChats.Next() {
		var chat entity.Chat
		if errChat := rowsChats.Scan(&chat.Id, &chat.Name, &chat.CreatedAt, &chat.IsDeleted); errChat != nil {
			return nil, fmt.Errorf("error path: %s, error: %w", opGetChat, errChat)
		}
		chats = append(chats, chat)
	}

	// В конце проверяем строки на ошибки (best practice)
	if err = rowsChats.Err(); err != nil {
		return nil, fmt.Errorf("error path: %s, error: %w", opGetChat, err)
	}

	return chats, nil
}
