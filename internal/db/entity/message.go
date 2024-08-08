package entity

// Message - сущность для работы с сообщениями
type Message struct {
	Id        int64  `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
	IsDeleted bool   `json:"isDeleted"`
}

// MessageAdd - сущность для отправки сообщения в чат от лица пользователя
type MessageAdd struct {
	ChatID int64  `json:"chatID"`
	UserID int64  `json:"userID"`
	Text   string `json:"text"`
}

// MessageUpdate - сущность для редактирования сообщения от лица пользователя
type MessageUpdate struct {
	MessageID int64  `json:"messageID"`
	UserID    int64  `json:"userID"`
	NewText   string `json:"newText"`
}
