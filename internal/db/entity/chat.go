package entity

// Chat - сущность для работы с чатом
type Chat struct {
	Id        int64  `json:"id" db:"id"`
	Name      string `json:"name" db:"name"`
	CreatedAt string `json:"created_at" db:"created_at"`
	IsDeleted bool   `json:"is_deleted" db:"is_deleted"`
}

// ChatAdd - сущность для создания чата между пользователями в бд
type ChatAdd struct {
	ChatName string  `json:"chatName"`
	Users    []int64 `json:"users"`
}

// ChatGet - сущность для получения чата пользователя из бд
type ChatGet struct {
	UserID int64 `json:"user_id"`
}
