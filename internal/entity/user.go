package entity

// User - сущность для работы с пользователями
type User struct {
	Id        int64  `json:"id"`
	Username  string `json:"username"`
	Password  string `json:"password"`
	CreatedAt string `json:"createdAt"`
	IsDeleted bool   `json:"isDeleted"`
}
