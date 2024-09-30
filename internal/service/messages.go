package service

import (
	"errors"
	"service-chat/internal/db"
	"service-chat/internal/db/entity"
	"service-chat/internal/dto"
)

type MessageService struct {
	repo db.Message
}

func NewMessageService(repo db.Message) *MessageService {
	return &MessageService{repo: repo}
}

// AddMessage - отправка сообщения в чат от лица пользователя
func (ms *MessageService) AddMessage(in dto.MessageAdd) (int, error) {
	// Если запрос пустой
	if in.ChatID == 0 || in.UserID == 0 {
		return 0, errors.New("empty chat_id or user_id")
	} else if in.Text == "" {
		return 0, errors.New("empty text")
	}

	dataDB := entity.MessageAdd{
		ChatID: in.ChatID,
		UserID: in.UserID,
		Text:   in.Text,
	}
	return ms.repo.AddMessage(dataDB)
}

// UpdateMessage - редактирование сообщения от лица пользователя
func (ms *MessageService) UpdateMessage(in dto.MessageUpdate) (int, error) {
	// Если запрос пустой
	if in.MessageID == 0 || in.UserID == 0 {
		return 0, errors.New("empty chat_id or user_id")
	} else if in.NewText == "" {
		return 0, errors.New("empty new_text")
	}

	dataDB := entity.MessageUpdate{
		MessageID: in.MessageID,
		UserID:    in.UserID,
		NewText:   in.NewText,
	}
	return ms.repo.UpdateMessage(dataDB)
}

// GetMessage - получение списка сообщений из конкретного чата
func (ms *MessageService) GetMessage(in dto.MessageGet, userID int) ([]entity.Message, error) {
	// Если пустой запрос
	if in.ChatID == 0 {
		return nil, errors.New("empty chat_id")
	} else if userID == 0 {
		return nil, errors.New("empty user_id")
	}

	dataDB := entity.MessageGet{
		ChatID: in.ChatID,
		Limit:  *in.Limit,
		Offset: *in.Offset,
		UserID: userID,
	}
	return ms.repo.GetMessage(dataDB)
}

// DeleteMessage - удаление сообщений
func (ms *MessageService) DeleteMessage(in dto.MessageDelete, userID int) ([]entity.DelMsg, error) {
	// Если запрос пустой
	if in.MessageIds == nil || len(*in.MessageIds) == 0 {
		return nil, errors.New("message_ids is empty")
	}
	if userID == 0 {
		return nil, errors.New("user_id is empty")
	}

	dataDB := entity.MessageDel{
		MsgIds: *in.MessageIds,
		UserID: userID,
	}
	return ms.repo.DeleteMessage(dataDB)
}
