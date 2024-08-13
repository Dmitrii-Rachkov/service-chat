package dto

// MessageAdd - структура запроса для ручки отправить сообщение от лица пользователя в чат
type MessageAdd struct {
	ChatID int64  `json:"chat_id" validate:"required"`
	UserID int64  `json:"user_id" validate:"required"`
	Text   string `json:"text" validate:"required"`
}
