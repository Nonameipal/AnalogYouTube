package dto

import "testing"

func TestMessageText(t *testing.T) {
	tests := []struct {
		name string
		req  CreateChatMessageRequest
		want string
	}{
		{name: "uses text first", req: CreateChatMessageRequest{Text: "hello", Content: "old"}, want: "hello"},
		{name: "falls back to content", req: CreateChatMessageRequest{Content: "old"}, want: "old"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.req.MessageText(); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}

func TestWebSocketMessageText(t *testing.T) {
	req := WebSocketChatMessageRequest{Text: "", Content: "from websocket"}
	if got := req.MessageText(); got != "from websocket" {
		t.Fatalf("expected fallback content, got %q", got)
	}
}

func TestCommentText(t *testing.T) {
	createReq := CreateCommentRequest{Text: "new", Content: "old"}
	if got := createReq.CommentText(); got != "new" {
		t.Fatalf("expected text field, got %q", got)
	}

	updateReq := UpdateCommentRequest{Content: "old"}
	if got := updateReq.CommentText(); got != "old" {
		t.Fatalf("expected content fallback, got %q", got)
	}
}

func TestDonationTargetUserID(t *testing.T) {
	if got := (CreateDonationRequest{ReceiverID: 2, RecipientID: 3}).TargetUserID(); got != 2 {
		t.Fatalf("expected receiver id, got %d", got)
	}
	if got := (CreateDonationRequest{RecipientID: 3}).TargetUserID(); got != 3 {
		t.Fatalf("expected recipient id fallback, got %d", got)
	}
}
