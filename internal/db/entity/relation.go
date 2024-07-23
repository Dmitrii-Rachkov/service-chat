package entity

// UsersChat - сущность для работы с зависимостями пользователей и чатов
type UsersChat struct {
	Id     int64 `json:"id"`
	UserID int64 `json:"userID"`
	ChatID int64 `json:"chatID"`
}

// ChatsMessages - сущность для работы с зависимостями чатов и сообщений
type ChatsMessages struct {
	Id          int64 `json:"id"`
	UsersChatID int64 `json:"usersChatID"`
	MessageID   int64 `json:"messageID"`
}
