package dbmodel

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
)

func TestToDomainModels(t *testing.T) {
	now := time.Now()
	categoryID := int64(9)
	videoID := int64(5)
	chatID := int64(7)

	category := Category{ID: 1, Name: "Go", Description: sql.NullString{String: "Backend", Valid: true}, CreatedAt: now, UpdatedAt: now}.ToDomain()
	if category.ID != 1 || category.Name != "Go" || category.Description != "Backend" {
		t.Fatalf("unexpected category: %+v", category)
	}

	chat := Chat{ID: 2, FirstUserID: 1, SecondUserID: 3, CreatedAt: now}.ToDomain()
	if chat.ID != 2 || chat.FirstUserID != 1 || chat.SecondUserID != 3 {
		t.Fatalf("unexpected chat: %+v", chat)
	}

	request := ChatRequest{ID: 3, SenderID: 1, ReceiverID: 2, Status: domain.ChatRequestStatusAccepted, ChatID: sql.NullInt64{Int64: chatID, Valid: true}, CreatedAt: now, UpdatedAt: now}.ToDomain()
	if request.ChatID == nil || *request.ChatID != int(chatID) {
		t.Fatalf("chat id was not copied: %+v", request.ChatID)
	}

	message := ChatMessage{ID: 4, ChatID: 2, SenderID: 1, Text: "hello", CreatedAt: now}.ToDomain()
	if message.Text != "hello" || message.ChatID != 2 {
		t.Fatalf("unexpected message: %+v", message)
	}

	comment := Comment{ID: 5, UserID: 1, VideoID: 2, Text: "nice", Status: domain.CommentStatusActive, CreatedAt: now, UpdatedAt: now}.ToDomain()
	if comment.Text != "nice" || comment.Status != domain.CommentStatusActive {
		t.Fatalf("unexpected comment: %+v", comment)
	}

	donation := Donation{ID: 6, SenderID: 1, ReceiverID: 2, VideoID: sql.NullInt64{Int64: videoID, Valid: true}, Amount: 10.5, Message: sql.NullString{String: "thanks", Valid: true}, CreatedAt: now}.ToDomain()
	if donation.VideoID == nil || *donation.VideoID != int(videoID) || donation.Message != "thanks" {
		t.Fatalf("unexpected donation: %+v", donation)
	}

	playlist := Playlist{ID: 7, Name: "Favorites", UserID: 1, Description: sql.NullString{String: "Best", Valid: true}, CreatedAt: now, UpdatedAt: now}.ToDomain()
	if playlist.Name != "Favorites" || playlist.UserID != 1 || playlist.Description != "Best" {
		t.Fatalf("unexpected playlist: %+v", playlist)
	}

	user := User{ID: 8, Username: "neo", Email: sql.NullString{String: "neo@gmail.com", Valid: true}, Password: "hash", Role: domain.UserRole, AvatarURL: sql.NullString{String: "/avatar.png", Valid: true}, Description: sql.NullString{String: "channel", Valid: true}, CreatedAt: now, UpdatedAt: now}.ToDomain()
	if user.Username != "neo" || user.Email != "neo@gmail.com" || user.AvatarURL != "/avatar.png" || user.Description != "channel" {
		t.Fatalf("unexpected user: %+v", user)
	}

	video := Video{ID: 9, AuthorID: 1, CategoryID: sql.NullInt64{Int64: categoryID, Valid: true}, Title: "video", Description: sql.NullString{String: "desc", Valid: true}, VideoURL: "/video.mp4", ThumbnailURL: sql.NullString{String: "/thumb.jpg", Valid: true}, Views: 10, Status: domain.VideoStatusActive, CreatedAt: now, UpdatedAt: now}.ToDomain()
	if video.CategoryID == nil || *video.CategoryID != int(categoryID) || video.Title != "video" || video.ThumbnailURL != "/thumb.jpg" {
		t.Fatalf("unexpected video: %+v", video)
	}

	quality := VideoQuality{ID: 10, VideoID: 9, Quality: "720p", VideoURL: "/720.mp4", CreatedAt: now}.ToDomain()
	if quality.VideoID != 9 || quality.Quality != "720p" || quality.VideoURL != "/720.mp4" {
		t.Fatalf("unexpected quality: %+v", quality)
	}
}

func TestToDomainNullPointers(t *testing.T) {
	if donation := (Donation{VideoID: sql.NullInt64{Valid: false}}).ToDomain(); donation.VideoID != nil {
		t.Fatalf("expected nil video id, got %+v", donation.VideoID)
	}

	if request := (ChatRequest{ChatID: sql.NullInt64{Valid: false}}).ToDomain(); request.ChatID != nil {
		t.Fatalf("expected nil chat id, got %+v", request.ChatID)
	}

	if video := (Video{CategoryID: sql.NullInt64{Valid: false}}).ToDomain(); video.CategoryID != nil {
		t.Fatalf("expected nil category id, got %+v", video.CategoryID)
	}
}
