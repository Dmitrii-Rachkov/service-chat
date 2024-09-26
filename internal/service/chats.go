package service

import (
	"errors"

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
	// Если запрос пустой
	if in.ChatName == "" || len(in.Users) == 0 {
		return 0, errors.New("chat_name or users is empty")
	}

	dataDB := entity.ChatAdd{
		ChatName: in.ChatName,
		Users:    in.Users,
	}
	return s.repo.CreateChat(dataDB)
}

// GetChat - получаем список чатов пользователя
func (s *ChatService) GetChat(in dto.ChatGet) ([]entity.Chat, error) {
	// Если запрос пустой
	if in.UserID == nil || *in.UserID == 0 {
		return nil, errors.New("user_id is empty")
	}

	dataDB := entity.ChatGet{
		UserID: *in.UserID,
	}
	return s.repo.GetChat(dataDB)
}

// DeleteChat - удаляем чаты
func (s *ChatService) DeleteChat(in dto.ChatDelete, userID int) ([]entity.DeletedChats, error) {
	// Если запрос пустой
	if in.ChatIds == nil || len(*in.ChatIds) == 0 {
		return nil, errors.New("chat_ids is empty")
	}
	if userID == 0 {
		return nil, errors.New("user_id is empty")
	}

	dataDB := entity.ChatDelete{
		ChatIds: *in.ChatIds,
		UserID:  userID,
	}
	return s.repo.DeleteChat(dataDB)
}
