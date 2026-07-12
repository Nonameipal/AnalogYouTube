package dto

type CreateCommentRequest struct {
	Text     string `json:"text"`
	Content  string `json:"content"`
	ParentID *int   `json:"parent_id,omitempty"`
}

type UpdateCommentRequest struct {
	Text    string `json:"text"`
	Content string `json:"content"`
}

func (r CreateCommentRequest) CommentText() string {
	if r.Text != "" {
		return r.Text
	}

	return r.Content
}

func (r UpdateCommentRequest) CommentText() string {
	if r.Text != "" {
		return r.Text
	}

	return r.Content
}
