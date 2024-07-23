package entity

// Chat - сущность для работы с чатом
type Chat struct {
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"createdAt"`
	IsDeleted bool   `json:"isDeleted"`
}
