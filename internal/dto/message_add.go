package dto

// MessageAdd - структура запроса для ручки отправить сообщение от лица пользователя в чат
type MessageAdd struct {
	ChatID int64  `json:"chat_id" validate:"required"`
	UserID int64  `json:"user_id" validate:"required"`
	Text   string `json:"text" validate:"required"`
}

// MessageUpdate - структура запроса для ручки редактирования сообщения от лица пользователя
type MessageUpdate struct {
	MessageID int64  `json:"message_id" validate:"required"`
	UserID    int64  `json:"user_id" validate:"required"`
	NewText   string `json:"new_text" validate:"required"`
}
