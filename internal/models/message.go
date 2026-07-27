package models

type Message struct {
    ID     string `json:"id"`
    UserID string `json:"user"`
    Text   string `json:"text"`
}
