package dto

// ChatDelete - структура запроса для ручки удаления чата
type ChatDelete struct {
	ChatIds *[]int64 `json:"chat_ids" validate:"required"`
}
