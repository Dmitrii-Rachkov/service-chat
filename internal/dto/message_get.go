package dto

// MessageGet - структура запроса для ручки получения списка сообщений в конкретном чате
type MessageGet struct {
	ChatID int64  `json:"chat_id" validate:"required"`
	Limit  *int64 `json:"limit" validate:"required"`
	Offset *int64 `json:"offset" validate:"required"`
}
