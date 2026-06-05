package dto

type CreateCommentRequest struct {
	Text string `json:"text"`
}

type UpdateCommentRequest struct {
	Text string `json:"text"`
}
