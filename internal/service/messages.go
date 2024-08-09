package service

import (
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
	dataDB := entity.MessageAdd{
		ChatID: in.ChatID,
		UserID: in.UserID,
		Text:   in.Text,
	}
	return ms.repo.AddMessage(dataDB)
}

// UpdateMessage - редактирование сообщения от лица пользователя
func (ms *MessageService) UpdateMessage(in dto.MessageUpdate) (int, error) {
	dataDB := entity.MessageUpdate{
		MessageID: in.MessageID,
		UserID:    in.UserID,
		NewText:   in.NewText,
	}
	return ms.repo.UpdateMessage(dataDB)
}
