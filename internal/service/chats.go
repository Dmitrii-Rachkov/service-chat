package service

import (
	"service-chat/internal/db"
	"service-chat/internal/db/entity"
	"service-chat/internal/dto"
)

type ChatService struct {
	repo db.Chat
}

func NewChatService(repo db.Chat) *ChatService {
	return &ChatService{repo: repo}
}

// CreateChat - создаём чат между пользователями
func (s *ChatService) CreateChat(in dto.ChatAdd) (int, error) {
	dataDB := entity.ChatAdd{
		ChatName: in.ChatName,
		Users:    in.Users,
	}
	return s.repo.CreateChat(dataDB)
}

// GetChat - получаем список чатов пользователя
func (s *ChatService) GetChat(in dto.ChatGet) ([]entity.Chat, error) {
	dataDB := entity.ChatGet{
		UserID: *in.UserID,
	}
	return s.repo.GetChat(dataDB)
}

// DeleteChat - удаляем чаты
func (s *ChatService) DeleteChat(in dto.ChatDelete, userID int) ([]entity.DeletedChats, error) {
	dataDB := entity.ChatDelete{
		ChatIds: *in.ChatIds,
		UserID:  userID,
	}
	return s.repo.DeleteChat(dataDB)
}
