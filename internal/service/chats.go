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
