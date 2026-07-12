package usecase

import (
	"errors"
	"testing"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func TestChatUsecase(t *testing.T) {
	t.Run("send request validates users and status", func(t *testing.T) {
		repo := newFakeRepository()
		repo.usersByID[1] = domain.User{ID: 1}
		repo.usersByID[2] = domain.User{ID: 2}
		uc := newTestUsecase(repo)

		request, err := uc.SendChatRequest(1, 2)
		if err != nil {
			t.Fatalf("SendChatRequest returned error: %v", err)
		}
		if request.SenderID != 1 || request.ReceiverID != 2 || request.Status != domain.ChatRequestStatusPending {
			t.Fatalf("unexpected request: %+v", request)
		}
	})

	t.Run("send request rejects existing chat", func(t *testing.T) {
		repo := newFakeRepository()
		repo.usersByID[1] = domain.User{ID: 1}
		repo.usersByID[2] = domain.User{ID: 2}
		repo.chatsByID[10] = domain.Chat{ID: 10, FirstUserID: 1, SecondUserID: 2}
		uc := newTestUsecase(repo)

		_, err := uc.SendChatRequest(1, 2)
		if !errors.Is(err, errs.ErrChatAlreadyExists) {
			t.Fatalf("expected ErrChatAlreadyExists, got %v", err)
		}
	})

	t.Run("accept request only by receiver", func(t *testing.T) {
		repo := newFakeRepository()
		repo.chatRequests[4] = domain.ChatRequest{ID: 4, SenderID: 1, ReceiverID: 2, Status: domain.ChatRequestStatusPending}
		uc := newTestUsecase(repo)

		accepted, err := uc.AcceptChatRequest(2, 4)
		if err != nil {
			t.Fatalf("AcceptChatRequest returned error: %v", err)
		}
		if accepted.Request.Status != domain.ChatRequestStatusAccepted || accepted.Chat.FirstUserID != 1 || accepted.Chat.SecondUserID != 2 {
			t.Fatalf("unexpected accepted request: %+v", accepted)
		}

		_, err = uc.AcceptChatRequest(1, 4)
		if !errors.Is(err, errs.ErrAccessDenied) {
			t.Fatalf("expected ErrAccessDenied for wrong receiver, got %v", err)
		}
	})

	t.Run("create message checks access", func(t *testing.T) {
		repo := newFakeRepository()
		repo.chatsByID[8] = domain.Chat{ID: 8, FirstUserID: 1, SecondUserID: 2}
		uc := newTestUsecase(repo)

		message, err := uc.CreateChatMessage(1, 8, domain.ChatMessage{Text: " hello "})
		if err != nil {
			t.Fatalf("CreateChatMessage returned error: %v", err)
		}
		if message.SenderID != 1 || message.ChatID != 8 || message.Text != "hello" {
			t.Fatalf("unexpected message: %+v", message)
		}

		_, err = uc.CreateChatMessage(3, 8, domain.ChatMessage{Text: "hello"})
		if !errors.Is(err, errs.ErrAccessDenied) {
			t.Fatalf("expected ErrAccessDenied, got %v", err)
		}
	})
}
