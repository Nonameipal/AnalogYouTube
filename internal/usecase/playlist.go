package usecase

import (
	"errors"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func (uc *Usecase) CreatePlaylist(userID int, playlist domain.Playlist) (domain.Playlist, error) {
	playlist.Name = strings.TrimSpace(playlist.Name)
	playlist.Description = strings.TrimSpace(playlist.Description)

	if userID <= 0 || playlist.Name == "" {
		return domain.Playlist{}, errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetUserByID(userID); err != nil {
		return domain.Playlist{}, err
	}

	playlist.UserID = userID

	return uc.repository.CreatePlaylist(playlist)
}

func (uc *Usecase) GetPlaylistByID(id int) (domain.Playlist, error) {
	if id <= 0 {
		return domain.Playlist{}, errs.ErrInvalidFieldValue
	}

	playlist, err := uc.repository.GetPlaylistByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.Playlist{}, errs.ErrPlaylistNotFound
		}
		return domain.Playlist{}, err
	}

	videos, err := uc.repository.GetPlaylistVideos(id)
	if err != nil {
		return domain.Playlist{}, err
	}
	playlist.Videos = videos
	return playlist, nil
}

func (uc *Usecase) GetUserPlaylists(userID int) ([]domain.Playlist, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}
	if _, err := uc.GetUserByID(userID); err != nil {
		return nil, err
	}
	return uc.repository.GetUserPlaylists(userID)
}

func (uc *Usecase) ensureCanManagePlaylist(userID int, userRole string, playlistID int) (domain.Playlist, error) {
	if userID <= 0 || playlistID <= 0 {
		return domain.Playlist{}, errs.ErrInvalidFieldValue
	}

	playlist, err := uc.repository.GetPlaylistByID(playlistID)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.Playlist{}, errs.ErrPlaylistNotFound
		}
		return domain.Playlist{}, err
	}
	if playlist.UserID != userID && userRole != domain.AdminRole {
		return domain.Playlist{}, errs.ErrAccessDenied
	}
	return playlist, nil
}

func (uc *Usecase) UpdatePlaylist(userID int, userRole string, playlist domain.Playlist) (domain.Playlist, error) {
	playlist.Name = strings.TrimSpace(playlist.Name)
	playlist.Description = strings.TrimSpace(playlist.Description)

	if playlist.ID <= 0 || playlist.Name == "" {
		return domain.Playlist{}, errs.ErrInvalidFieldValue
	}
	if _, err := uc.ensureCanManagePlaylist(userID, userRole, playlist.ID); err != nil {
		return domain.Playlist{}, err
	}
	return uc.repository.UpdatePlaylist(playlist)
}

func (uc *Usecase) DeletePlaylist(userID int, userRole string, playlistID int) error {
	if _, err := uc.ensureCanManagePlaylist(userID, userRole, playlistID); err != nil {
		return err
	}
	if err := uc.repository.DeletePlaylist(playlistID); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrPlaylistNotFound
		}
		return err
	}
	return nil
}

func (uc *Usecase) AddVideoToPlaylist(userID int, userRole string, playlistID int, videoID int) error {
	if videoID <= 0 {
		return errs.ErrInvalidFieldValue
	}

	if _, err := uc.ensureCanManagePlaylist(userID, userRole, playlistID); err != nil {
		return err
	}
	if _, err := uc.GetVideoByID(videoID); err != nil {
		return err
	}
	return uc.repository.AddVideoToPlaylist(playlistID, videoID)
}

func (uc *Usecase) RemoveVideoFromPlaylist(userID int, userRole string, playlistID int, videoID int) error {
	if videoID <= 0 {
		return errs.ErrInvalidFieldValue
	}
	if _, err := uc.ensureCanManagePlaylist(userID, userRole, playlistID); err != nil {
		return err
	}
	return uc.repository.RemoveVideoFromPlaylist(playlistID, videoID)
}
