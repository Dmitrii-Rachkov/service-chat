package dto

// ChatGet - структура запроса для ручки получения списка чатов конкретного пользователя
type ChatGet struct {
	UserID *int64 `json:"user_id" validate:"required"`
}
