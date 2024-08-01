package dto

// ChatAdd - структура запроса для ручки создания чата
type ChatAdd struct {
	ChatName string  `json:"chat_name" validate:"required,max=20,min=6,excludesall=!@#$&*()?"`
	Users    []int64 `json:"users" validate:"required"`
}
