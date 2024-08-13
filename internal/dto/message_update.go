package dto

// MessageUpdate - структура запроса для ручки редактирования сообщения от лица пользователя
type MessageUpdate struct {
	MessageID int64  `json:"message_id" validate:"required"`
	UserID    int64  `json:"user_id" validate:"required"`
	NewText   string `json:"new_text" validate:"required"`
}
