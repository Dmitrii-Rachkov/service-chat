package dto

// MessageDelete - структура запроса для ручки удаления сообщение от лица пользователя
type MessageDelete struct {
	MessageIds *[]int64 `json:"message_ids" validate:"required"`
}
