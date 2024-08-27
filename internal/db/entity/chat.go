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

// ChatDelete - сущность для soft удаления чатов в бд
type ChatDelete struct {
	ChatIds []int64 `json:"chat_ids"`
	UserID  int     `json:"user_id"`
}

// DeletedChats - сущность для получения результата удаления чатов из бд
type DeletedChats struct {
	ChatID int64  `json:"chat_id" db:"identifier"`
	Result string `json:"result" db:"result"`
}
