package dto

type CreateChatRequest struct {
	UserID int `json:"user_id"`
}

type CreateChatMessageRequest struct {
	Text    string `json:"text"`
	Content string `json:"content"`
}

type WebSocketChatMessageRequest struct {
	Text    string `json:"text"`
	Content string `json:"content"`
}

func (r CreateChatMessageRequest) MessageText() string {
	if r.Text != "" {
		return r.Text
	}

	return r.Content
}

func (r WebSocketChatMessageRequest) MessageText() string {
	if r.Text != "" {
		return r.Text
	}

	return r.Content
}
