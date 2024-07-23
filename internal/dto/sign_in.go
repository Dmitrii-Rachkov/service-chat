package dto

// SignInRequest - структура запроса для ручки авторизации пользователя
type SignInRequest struct {
	Username string `json:"username" validate:"required,max=20,excludesall=!@#$&*()?"`
	Password string `json:"password" validate:"required,max=12,min=6,containsany=@#$&*()"`
}
