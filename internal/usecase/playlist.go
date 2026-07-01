package usecase

import (
	"errors"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func (s *Service) CreatePlaylist(userID int, playlist domain.Playlist) (domain.Playlist, error) {
	playlist.Name = strings.TrimSpace(playlist.Name)
	playlist.Description = strings.TrimSpace(playlist.Description)

	if userID <= 0 || playlist.Name == "" {
		return domain.Playlist{}, errs.ErrInvalidFieldValue
	}

	if _, err := s.GetUserByID(userID); err != nil {
		return domain.Playlist{}, err
	}

	playlist.UserID = userID

	return s.repository.CreatePlaylist(playlist)
}

func (s *Service) GetPlaylistByID(id int) (domain.Playlist, error) {
	if id <= 0 {
		return domain.Playlist{}, errs.ErrInvalidFieldValue
	}

	playlist, err := s.repository.GetPlaylistByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.Playlist{}, errs.ErrPlaylistNotFound
		}
		return domain.Playlist{}, err
	}

	videos, err := s.repository.GetPlaylistVideos(id)
	if err != nil {
		return domain.Playlist{}, err
	}
	playlist.Videos = videos
	return playlist, nil
}

func (s *Service) GetUserPlaylists(userID int) ([]domain.Playlist, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}
	if _, err := s.GetUserByID(userID); err != nil {
		return nil, err
	}
	return s.repository.GetUserPlaylists(userID)
}

func (s *Service) ensureCanManagePlaylist(userID int, userRole string, playlistID int) (domain.Playlist, error) {
	if userID <= 0 || playlistID <= 0 {
		return domain.Playlist{}, errs.ErrInvalidFieldValue
	}

	playlist, err := s.repository.GetPlaylistByID(playlistID)
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

func (s *Service) UpdatePlaylist(userID int, userRole string, playlist domain.Playlist) (domain.Playlist, error) {
	playlist.Name = strings.TrimSpace(playlist.Name)
	playlist.Description = strings.TrimSpace(playlist.Description)

	if playlist.ID <= 0 || playlist.Name == "" {
		return domain.Playlist{}, errs.ErrInvalidFieldValue
	}
	if _, err := s.ensureCanManagePlaylist(userID, userRole, playlist.ID); err != nil {
		return domain.Playlist{}, err
	}
	return s.repository.UpdatePlaylist(playlist)
}

func (s *Service) DeletePlaylist(userID int, userRole string, playlistID int) error {
	if _, err := s.ensureCanManagePlaylist(userID, userRole, playlistID); err != nil {
		return err
	}
	if err := s.repository.DeletePlaylist(playlistID); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrPlaylistNotFound
		}
		return err
	}
	return nil
}

func (s *Service) AddVideoToPlaylist(userID int, userRole string, playlistID int, videoID int) error {
	if videoID <= 0 {
		return errs.ErrInvalidFieldValue
	}

	if _, err := s.ensureCanManagePlaylist(userID, userRole, playlistID); err != nil {
		return err
	}
	if _, err := s.GetVideoByID(videoID); err != nil {
		return err
	}
	return s.repository.AddVideoToPlaylist(playlistID, videoID)
}

func (s *Service) RemoveVideoFromPlaylist(userID int, userRole string, playlistID int, videoID int) error {
	if videoID <= 0 {
		return errs.ErrInvalidFieldValue
	}
	if _, err := s.ensureCanManagePlaylist(userID, userRole, playlistID); err != nil {
		return err
	}
	return s.repository.RemoveVideoFromPlaylist(playlistID, videoID)
}
