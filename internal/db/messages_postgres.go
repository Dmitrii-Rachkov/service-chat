package db

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/lib/pq"

	"service-chat/internal/db/entity"
)

const (
	opMessageAdd    = "db.AddMessage"
	opMessageUpdate = "db.UpdateMessage"
	opMessageGet    = "db.GetMessage"
	opDelMsg        = "db.DeleteMessage"
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

// GetMessage - получаем список сообщений в конкретном чате из бд
func (m *MessagePostgres) GetMessage(in entity.MessageGet) ([]entity.Message, error) {
	// Пользователь может достать сообщения из чата, только если он состоит в чате,
	// для этого нужно SELECT из таблицы users_chat, где chat_id = in.ChatID,
	// получаем users_chat_id (несколько) и user_id (несколько), которые в этом чате состоят

	// Скелет sql запроса на получение users_chat_id и user_id из конкретного чата
	stmtUsersChat, err := m.db.Prepare(`SELECT id, user_id FROM "users_chat" WHERE chat_id = $1`)
	if err != nil {
		return nil, fmt.Errorf("error path: %s, error: %w", opMessageGet, err)
	}
	defer stmtUsersChat.Close()

	// Получаем users_chat_id и user_id из бд
	rowsUc, err := stmtUsersChat.Query(in.ChatID)
	if err != nil {
		return nil, fmt.Errorf("error path: %s, error: %w", opMessageGet, err)
	}
	defer rowsUc.Close()

	// Кладём строки из бд в структуру
	type usersChat struct {
		usersChatID []int `db:"id"`
		userID      []int `db:"user_id"`
	}
	var uc usersChat
	for rowsUc.Next() {
		var usersChatID, userID int
		if errSc := rowsUc.Scan(&usersChatID, &userID); errSc != nil {
			return nil, fmt.Errorf("error path: %s, error: %w", opMessageGet, errSc)
		}
		uc.usersChatID = append(uc.usersChatID, usersChatID)
		uc.userID = append(uc.userID, userID)
	}

	// В конце проверяем строки на ошибки (best practice)
	if err = rowsUc.Err(); err != nil {
		return nil, fmt.Errorf("error path: %s, error: %w", opMessageGet, err)
	}

	// Проверяем, что пользователь, который запрашивает сообщения, есть в чате
	// Пользователь не может запрашивать сообщения из чата, если он в нём не состоит
	var checkUser bool
	for _, id := range uc.userID {
		if id == in.UserID {
			checkUser = true
			break
		}
	}
	if !checkUser {
		return nil, fmt.Errorf("error path: %s, error: %w", opMessageGet,
			errors.New(fmt.Sprintf("User with userID %d does not exist in chatID %d", in.UserID, in.ChatID)))
	}

	// Далее по users_chat_id (несколько) берём все messageID из таблицы chats_messages,
	// и забираем все сообщения из message

	// Скелет sql запроса на получение всех сообщений в конкретном чате
	stmtMsg, err := m.db.Prepare(`WITH cm AS (
											SELECT message_id
											FROM chats_messages
											WHERE users_chat_id = ANY ($1)
											)
										SELECT *
										FROM message
										WHERE id IN (SELECT message_id FROM cm)
										ORDER BY created_at
										LIMIT $2 OFFSET $3`)
	if err != nil {
		return nil, fmt.Errorf("error path: %s, error: %w", opMessageGet, err)
	}
	defer stmtMsg.Close()

	// Получаем сообщения из бд
	rowsMsg, err := stmtMsg.Query(pq.Array(uc.usersChatID), in.Limit, in.Offset)
	if err != nil {
		return nil, fmt.Errorf("error path: %s, error: %w", opMessageGet, err)
	}
	defer rowsMsg.Close()

	// Структура для записи всех полученных сообщений из бд
	var messages []entity.Message
	for rowsMsg.Next() {
		var msg entity.Message
		if errSc := rowsMsg.Scan(&msg.Id, &msg.Text, &msg.UserID, &msg.CreatedAt, &msg.IsDeleted); errSc != nil {
			return nil, fmt.Errorf("error path: %s, error: %w", opMessageGet, errSc)
		}
		messages = append(messages, msg)
	}

	// В конце проверяем строки на ошибки (best practice)
	if err = rowsMsg.Err(); err != nil {
		return nil, fmt.Errorf("error path: %s, error: %w", opMessageGet, err)
	}

	return messages, nil
}

// DeleteMessage - soft удаление сообщений
func (m *MessagePostgres) DeleteMessage(in entity.MessageDel) ([]entity.DelMsg, error) {
	// Скелет sql запроса на удаление сообщений
	stmtDelMsg, err := m.db.Prepare(`SELECT * FROM delete_message($1, VARIADIC $2::integer[])`)
	if err != nil {
		return nil, fmt.Errorf("error path: %s, error: %w", opDelMsg, err)
	}
	defer stmtDelMsg.Close()

	// Получаем информацию по результатам soft удаления сообщений из бд
	rowsDel, err := stmtDelMsg.Query(in.UserID, pq.Array(in.MsgIds))
	if err != nil {
		return nil, fmt.Errorf("error path: %s, error: %w", opDelMsg, err)
	}
	defer rowsDel.Close()

	// Структура для записи результата soft удаления сообщений из бд
	var delMsg []entity.DelMsg
	for rowsDel.Next() {
		var dm entity.DelMsg
		if errDel := rowsDel.Scan(&dm.MessageID, &dm.Result); errDel != nil {
			return nil, fmt.Errorf("error path: %s, error: %w", opDelMsg, errDel)
		}
		delMsg = append(delMsg, dm)
	}

	// В конце проверяем строки на ошибки (best practice)
	if err = rowsDel.Err(); err != nil {
		return nil, fmt.Errorf("error path: %s, error: %w", opDelMsg, err)
	}

	return delMsg, nil
}
