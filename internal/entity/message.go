package entity

// Message - сущность для работы с сообщениями
type Message struct {
	Id        int64  `json:"id"`
	Text      string `json:"text"`
	CreatedAt string `json:"createdAt"`
	IsDeleted bool   `json:"isDeleted"`
}
