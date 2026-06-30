package dto

type CreateChatRequest struct {
	UserID int `json:"user_id"`
}

type CreateChatMessageRequest struct {
	Text string `json:"text"`
}

type WebSocketChatMessageRequest struct {
	Text string `json:"text"`
}
