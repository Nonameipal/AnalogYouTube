package usecase

import (
	"errors"
	"reflect"
	"testing"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func TestVideoUsecase(t *testing.T) {
	t.Run("create allows empty category and saves active video", func(t *testing.T) {
		repo := newFakeRepository()
		uc := newTestUsecase(repo)

		video, err := uc.CreateVideo(1, domain.Video{Title: " My video ", VideoURL: " /video.mp4 ", ThumbnailURL: " /thumb.jpg "})
		if err != nil {
			t.Fatalf("CreateVideo returned error: %v", err)
		}
		if video.AuthorID != 1 || video.Status != domain.VideoStatusActive || video.CategoryID != nil {
			t.Fatalf("unexpected created video: %+v", video)
		}
		if video.Title != "My video" || video.VideoURL != "/video.mp4" {
			t.Fatalf("video fields were not trimmed: %+v", video)
		}
	})

	t.Run("create resolves category by name", func(t *testing.T) {
		repo := newFakeRepository()
		repo.categories["Go"] = domain.Category{ID: 9, Name: "Go"}
		uc := newTestUsecase(repo)
		categoryName := " Go "

		video, err := uc.CreateVideo(1, domain.Video{Title: "video", VideoURL: "/video.mp4", CategoryName: &categoryName})
		if err != nil {
			t.Fatalf("CreateVideo returned error: %v", err)
		}
		if video.CategoryID == nil || *video.CategoryID != 9 {
			t.Fatalf("expected category id 9, got %+v", video.CategoryID)
		}
	})

	t.Run("get adds qualities and maps not found", func(t *testing.T) {
		repo := newFakeRepository()
		repo.videosByID[1] = domain.Video{ID: 1, Title: "video"}
		repo.videoQualities[1] = []domain.VideoQuality{{Quality: "720p"}}
		uc := newTestUsecase(repo)

		video, err := uc.GetVideoByID(1)
		if err != nil {
			t.Fatalf("GetVideoByID returned error: %v", err)
		}
		if len(video.Qualities) != 1 || video.Qualities[0].Quality != "720p" {
			t.Fatalf("qualities were not attached: %+v", video.Qualities)
		}

		_, err = uc.GetVideoByID(404)
		if !errors.Is(err, errs.ErrVideoNotFound) {
			t.Fatalf("expected ErrVideoNotFound, got %v", err)
		}
	})

	t.Run("update keeps old file and thumbnail", func(t *testing.T) {
		repo := newFakeRepository()
		oldCategoryID := 3
		repo.videosByID[1] = domain.Video{ID: 1, AuthorID: 10, VideoURL: "/old.mp4", ThumbnailURL: "/old.jpg", CategoryID: &oldCategoryID}
		uc := newTestUsecase(repo)

		video, err := uc.UpdateVideo(10, domain.UserRole, domain.Video{ID: 1, Title: " new "})
		if err != nil {
			t.Fatalf("UpdateVideo returned error: %v", err)
		}
		if video.VideoURL != "/old.mp4" || video.ThumbnailURL != "/old.jpg" {
			t.Fatalf("old file fields were not kept: %+v", video)
		}
		if video.CategoryID == nil || *video.CategoryID != 3 {
			t.Fatalf("old category was not kept: %+v", video.CategoryID)
		}
	})

	t.Run("update denies another user", func(t *testing.T) {
		repo := newFakeRepository()
		repo.videosByID[1] = domain.Video{ID: 1, AuthorID: 10, Title: "video"}
		uc := newTestUsecase(repo)

		_, err := uc.UpdateVideo(99, domain.UserRole, domain.Video{ID: 1, Title: "new"})
		if !errors.Is(err, errs.ErrAccessDenied) {
			t.Fatalf("expected ErrAccessDenied, got %v", err)
		}
	})

	t.Run("delete allows admin", func(t *testing.T) {
		repo := newFakeRepository()
		repo.videosByID[1] = domain.Video{ID: 1, AuthorID: 10}
		uc := newTestUsecase(repo)

		err := uc.DeleteVideo(99, domain.AdminRole, 1)
		if err != nil {
			t.Fatalf("DeleteVideo returned error: %v", err)
		}
		if repo.deletedVideoID != 1 {
			t.Fatalf("expected deleted video id 1, got %d", repo.deletedVideoID)
		}
	})

	t.Run("search empty title returns recommended", func(t *testing.T) {
		repo := newFakeRepository()
		repo.recommendedVideos = []domain.Video{{ID: 1, Title: "recommended"}}
		uc := newTestUsecase(repo)

		videos, err := uc.SearchVideosByTitle("   ")
		if err != nil {
			t.Fatalf("SearchVideosByTitle returned error: %v", err)
		}
		if !reflect.DeepEqual(videos, repo.recommendedVideos) {
			t.Fatalf("expected recommended videos, got %+v", videos)
		}
	})

	t.Run("playback speeds are stable", func(t *testing.T) {
		got := newTestUsecase(newFakeRepository()).GetPlaybackSpeeds()
		want := []float64{0.25, 1.0, 1.25, 1.5, 2.0}
		if !reflect.DeepEqual(got, want) {
			t.Fatalf("expected speeds %+v, got %+v", want, got)
		}
	})
}
