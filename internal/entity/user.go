package entity

// User - сущность для работы с пользователями
type User struct {
	Id        int64  `json:"id"`
	Username  string `json:"username" validate:"required,max=20"`
	Password  string `json:"password" validate:"required,max=12,min=6,containsany=@#$&*()"`
	CreatedAt string `json:"createdAt"`
	IsDeleted bool   `json:"isDeleted"`
}
