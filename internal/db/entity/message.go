package entity

// Message - сущность для работы с сообщениями
type Message struct {
	Id        int64  `json:"id" db:"id"`
	Text      string `json:"text" db:"text"`
	UserID    int64  `json:"user_id" db:"user_id"`
	CreatedAt string `json:"created_at" db:"created_at"`
	IsDeleted bool   `json:"is_deleted" db:"is_deleted"`
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

// MessageGet - сущность для получения списка сообщений в конкретном чате
type MessageGet struct {
	ChatID int64 `json:"chatID"`
	Limit  int64 `json:"limit"`
	Offset int64 `json:"offset"`
	UserID int   `json:"userID"`
}

// MessageDel - сущность для удаления сообщений
type MessageDel struct {
	MsgIds []int64 `json:"message_Ids"`
	UserID int     `json:"user_id"`
}

// DelMsg - сущность для получения результата удаления сообщений из бд
type DelMsg struct {
	MessageID int64  `json:"message_id" db:"identifier"`
	Result    string `json:"result" db:"result"`
}
