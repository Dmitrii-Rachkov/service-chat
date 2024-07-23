package dto

// SignUpRequest - структура запроса для ручки регистрации пользователя
type SignUpRequest struct {
	Username string `json:"username" validate:"required,max=20,excludesall=!@#$&*()?"`
	Password string `json:"password" validate:"required,max=12,min=6,containsany=@#$&*()"`
}
