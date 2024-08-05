package entity

// Chat - сущность для работы с чатом
type Chat struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	IsDeleted bool   `json:"isDeleted"`
}

// ChatAdd - сущность для создания чата между пользователями в бд
type ChatAdd struct {
	ChatName string  `json:"chatName"`
	Users    []int64 `json:"users"`
}
