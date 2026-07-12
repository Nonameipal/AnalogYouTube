package usecase

import (
	"errors"
	"testing"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func TestCommentUsecase(t *testing.T) {
	t.Run("create trims and links comment", func(t *testing.T) {
		repo := newFakeRepository()
		repo.videosByID[5] = domain.Video{ID: 5, AuthorID: 2}
		uc := newTestUsecase(repo)

		comment, err := uc.CreateComment(1, 5, domain.Comment{Text: " nice "})
		if err != nil {
			t.Fatalf("CreateComment returned error: %v", err)
		}
		if comment.UserID != 1 || comment.VideoID != 5 || comment.Text != "nice" || comment.Status != domain.CommentStatusActive {
			t.Fatalf("unexpected comment: %+v", comment)
		}
		if comment.ParentID != nil {
			t.Fatalf("top-level comment should have nil ParentID, got %v", comment.ParentID)
		}
	})

	t.Run("reply to comment saves parent_id", func(t *testing.T) {
		repo := newFakeRepository()
		repo.videosByID[5] = domain.Video{ID: 5, AuthorID: 2}
		repo.commentsByID[10] = domain.Comment{ID: 10, UserID: 2, VideoID: 5, Text: "original"}
		uc := newTestUsecase(repo)

		parentID := 10
		reply, err := uc.CreateComment(1, 5, domain.Comment{Text: "reply!", ParentID: &parentID})
		if err != nil {
			t.Fatalf("CreateComment (reply) returned error: %v", err)
		}
		if reply.ParentID == nil || *reply.ParentID != 10 {
			t.Fatalf("expected ParentID=10, got %v", reply.ParentID)
		}
	})

	t.Run("reply to comment on different video is rejected", func(t *testing.T) {
		repo := newFakeRepository()
		repo.videosByID[5] = domain.Video{ID: 5, AuthorID: 2}
		repo.videosByID[9] = domain.Video{ID: 9, AuthorID: 3}
		repo.commentsByID[10] = domain.Comment{ID: 10, UserID: 2, VideoID: 9, Text: "other video comment"}
		uc := newTestUsecase(repo)

		parentID := 10
		_, err := uc.CreateComment(1, 5, domain.Comment{Text: "sneaky reply", ParentID: &parentID})
		if !errors.Is(err, errs.ErrInvalidFieldValue) {
			t.Fatalf("expected ErrInvalidFieldValue, got %v", err)
		}
	})

	t.Run("reply to non-existent comment is rejected", func(t *testing.T) {
		repo := newFakeRepository()
		repo.videosByID[5] = domain.Video{ID: 5, AuthorID: 2}
		uc := newTestUsecase(repo)

		parentID := 999
		_, err := uc.CreateComment(1, 5, domain.Comment{Text: "reply to ghost", ParentID: &parentID})
		if !errors.Is(err, errs.ErrCommentNotFound) {
			t.Fatalf("expected ErrCommentNotFound, got %v", err)
		}
	})

	t.Run("update denies non owner", func(t *testing.T) {
		repo := newFakeRepository()
		repo.commentsByID[1] = domain.Comment{ID: 1, UserID: 10, Text: "old"}
		uc := newTestUsecase(repo)

		_, err := uc.UpdateComment(99, domain.UserRole, domain.Comment{ID: 1, Text: "new"})
		if !errors.Is(err, errs.ErrAccessDenied) {
			t.Fatalf("expected ErrAccessDenied, got %v", err)
		}
	})

	t.Run("admin can delete", func(t *testing.T) {
		repo := newFakeRepository()
		repo.commentsByID[1] = domain.Comment{ID: 1, UserID: 10}
		uc := newTestUsecase(repo)

		err := uc.DeleteComment(99, domain.AdminRole, 1)
		if err != nil {
			t.Fatalf("DeleteComment returned error: %v", err)
		}
		if repo.deletedCommentID != 1 {
			t.Fatalf("expected deleted comment id 1, got %d", repo.deletedCommentID)
		}
	})
}
