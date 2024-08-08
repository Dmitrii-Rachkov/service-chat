package db

import (
	"database/sql"
	"fmt"

	"service-chat/internal/db/entity"
)

const (
	opMessageAdd    = "db.AddMessage"
	opMessageUpdate = "db.UpdateMessage"
)

type MessagePostgres struct {
	db *sql.DB
}

func NewMessagePostgres(db *sql.DB) *MessagePostgres {
	return &MessagePostgres{db: db}
}

// AddMessage - сохраняем сообщение в чат от пользователя в бд и возвращаем message id
func (m *MessagePostgres) AddMessage(in entity.MessageAdd) (int, error) {
	var messageID int
	var cmID int

	// Начинаем транзакцию
	tx, err := m.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", opCreateChat, err)
	}

	// Скелет sql запроса на сохранение сообщения в бд
	stmtAdd, errAdd := tx.Prepare(`INSERT INTO "message" (text, user_id) VALUES ($1, $2) RETURNING id`)
	if errAdd != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", opMessageAdd, errAdd)
	}
	defer stmtAdd.Close()

	// Сохраняем сообщение от пользователя в бд
	if rowAdd := tx.Stmt(stmtAdd).QueryRow(in.Text, in.UserID).Scan(&messageID); rowAdd != nil {
		// Откатываем транзакцию в случае ошибки
		errTx := tx.Rollback()
		if errTx != nil {
			return 0, fmt.Errorf("error path: %s, error: %s", opCreateChat, errTx)
		}
		return 0, fmt.Errorf("error path: %s, error: %w", opMessageAdd, rowAdd)
	}

	// Скелет sql на связь users_chat_id и message_id в таблице chats_messages
	stmtCm, errCm := tx.Prepare(`WITH uci AS (
											SELECT id FROM "users_chat" 
											WHERE user_id = $1 AND chat_id = $2
										)
										INSERT INTO "chats_messages" (users_chat_id, message_id)
										SELECT id, $3 FROM uci 
										RETURNING id`)
	if errCm != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", opCreateChat, errCm)
	}
	defer stmtCm.Close()

	// Создаём связь users_chat_id и message_id в таблице chats_messages
	if rowCm := tx.Stmt(stmtCm).QueryRow(in.UserID, in.ChatID, messageID).Scan(&cmID); rowCm != nil && rowCm.Error() == errNoRows {
		errTx := tx.Rollback()
		if errTx != nil {
			return 0, fmt.Errorf("error path: %s, error: %s", opCreateChat, errTx)
		}
		return 0, fmt.Errorf("error path: %s, error: %s", opMessageAdd, "Invalid chat_id")
	} else if rowCm != nil {
		errTx := tx.Rollback()
		if errTx != nil {
			return 0, fmt.Errorf("error path: %s, error: %s", opCreateChat, errTx)
		}
		return 0, fmt.Errorf("error path: %s, error: %w", opMessageAdd, rowCm)
	}

	return messageID, tx.Commit()
}

// UpdateMessage - редактируем сообщение от пользователя в бд и возвращаем message id
func (m *MessagePostgres) UpdateMessage(in entity.MessageUpdate) (int, error) {
	var messageID int

	// Скелет sql запроса на редактирование сообщения в бд
	stmt, err := m.db.Prepare(`UPDATE "message" SET text = $1 WHERE id = $2 AND user_id = $3 RETURNING id`)
	if err != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", opMessageUpdate, err)
	}
	defer stmt.Close()

	// Редактируем сообщение от пользователя в бд
	if row := stmt.QueryRow(in.NewText, in.MessageID, in.UserID).Scan(&messageID); row != nil && row.Error() == errNoRows {
		return 0, fmt.Errorf("error path: %s, error: %s", opMessageUpdate, "Invalid message_id OR user_id")
	} else if row != nil {
		return 0, fmt.Errorf("error path: %s, error: %w", opMessageUpdate, row)
	}

	return messageID, nil
}
