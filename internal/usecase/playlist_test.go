package usecase

import (
	"errors"
	"testing"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func TestPlaylistUsecase(t *testing.T) {
	t.Run("create validates user and saves owner", func(t *testing.T) {
		repo := newFakeRepository()
		repo.usersByID[1] = domain.User{ID: 1}
		uc := newTestUsecase(repo)

		playlist, err := uc.CreatePlaylist(1, domain.Playlist{Name: " Favorites ", Description: " Best "})
		if err != nil {
			t.Fatalf("CreatePlaylist returned error: %v", err)
		}
		if playlist.UserID != 1 || playlist.Name != "Favorites" || playlist.Description != "Best" {
			t.Fatalf("unexpected playlist: %+v", playlist)
		}
	})

	t.Run("get attaches videos", func(t *testing.T) {
		repo := newFakeRepository()
		repo.playlistsByID[1] = domain.Playlist{ID: 1, UserID: 1, Name: "Favorites"}
		repo.playlistVideos[1] = []domain.Video{{ID: 2, Title: "video"}}
		uc := newTestUsecase(repo)

		playlist, err := uc.GetPlaylistByID(1)
		if err != nil {
			t.Fatalf("GetPlaylistByID returned error: %v", err)
		}
		if len(playlist.Videos) != 1 || playlist.Videos[0].ID != 2 {
			t.Fatalf("videos were not attached: %+v", playlist.Videos)
		}
	})

	t.Run("owner can add video", func(t *testing.T) {
		repo := newFakeRepository()
		repo.playlistsByID[1] = domain.Playlist{ID: 1, UserID: 1}
		repo.videosByID[2] = domain.Video{ID: 2}
		uc := newTestUsecase(repo)

		err := uc.AddVideoToPlaylist(1, domain.UserRole, 1, 2)
		if err != nil {
			t.Fatalf("AddVideoToPlaylist returned error: %v", err)
		}
		if repo.addedPlaylistID != 1 || repo.addedVideoID != 2 {
			t.Fatal("video was not added to playlist")
		}
	})

	t.Run("another user cannot update playlist", func(t *testing.T) {
		repo := newFakeRepository()
		repo.playlistsByID[1] = domain.Playlist{ID: 1, UserID: 1}
		uc := newTestUsecase(repo)

		_, err := uc.UpdatePlaylist(2, domain.UserRole, domain.Playlist{ID: 1, Name: "new"})
		if !errors.Is(err, errs.ErrAccessDenied) {
			t.Fatalf("expected ErrAccessDenied, got %v", err)
		}
	})
}
