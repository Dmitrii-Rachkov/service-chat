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
	stmtChats, err := c.db.Prepare(`WITH users_chat_ids AS (
												SELECT id
												FROM users_chat
												WHERE user_id = $1
											), chat_ids AS (
												SELECT chat_id
												FROM users_chat
												WHERE user_id = $1
											), with_messages AS (
												SELECT c.id, c.name, c.created_at, c.is_deleted
												FROM (
														 SELECT DISTINCT ON (users_chat_id) users_chat_id, message_id
														 FROM chats_messages
														 WHERE users_chat_id IN (SELECT id FROM users_chat_ids)
														 ORDER BY users_chat_id, message_id DESC
													 ) AS ls
												LEFT JOIN users_chat AS uc
												ON ls.users_chat_id = uc.id
												LEFT JOIN chat AS c
												ON uc.chat_id = c.id
												ORDER BY message_id DESC
											), no_message AS (
												SELECT *
												FROM chat
												WHERE id IN (SELECT chat_id FROM chat_ids)
												AND id NOT IN (SELECT id FROM with_messages)
												ORDER BY created_at DESC
											)
											SELECT *
											FROM with_messages
											UNION ALL
											SELECT *
											FROM no_message
`)

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
