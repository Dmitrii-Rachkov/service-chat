package entity

// User - сущность для работы с пользователями
type User struct {
	Id        int64  `json:"id" db:"id"`
	Username  string `json:"username" validate:"required,max=20" db:"username"`
	Password  string `json:"password" validate:"required,max=12,min=6,containsany=@#$&*()" db:"password_hash"`
	CreatedAt string `json:"createdAt"`
	IsDeleted bool   `json:"isDeleted"`
}
