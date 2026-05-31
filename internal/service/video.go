package service

import (
	"errors"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
)

func (s *Service) CreateVideo(authorID int, video domain.Video) (domain.Video, error) {
	video.Title = strings.TrimSpace(video.Title)
	video.Description = strings.TrimSpace(video.Description)
	video.VideoURL = strings.TrimSpace(video.VideoURL)
	video.ThumbnailURL = strings.TrimSpace(video.ThumbnailURL)

	if authorID <= 0 || video.Title == "" || video.VideoURL == "" {
		return domain.Video{}, errs.ErrInvalidFieldValue
	}

	video.AuthorID = authorID
	video.Status = domain.VideoStatusActive

	return s.repository.CreateVideo(video)
}

func (s *Service) GetAllVideos() ([]domain.Video, error) {
	return s.repository.GetAllVideos()
}

func (s *Service) GetVideoByID(id int) (domain.Video, error) {
	if id <= 0 {
		return domain.Video{}, errs.ErrInvalidFieldValue
	}

	video, err := s.repository.GetVideoByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.Video{}, errs.ErrVideoNotFound
		}
		return domain.Video{}, err
	}

	return video, nil
}

func (s *Service) UpdateVideo(userID int, userRole string, video domain.Video) (domain.Video, error) {
	video.Title = strings.TrimSpace(video.Title)
	video.Description = strings.TrimSpace(video.Description)
	video.VideoURL = strings.TrimSpace(video.VideoURL)
	video.ThumbnailURL = strings.TrimSpace(video.ThumbnailURL)

	if userID <= 0 || video.ID <= 0 || video.Title == "" || video.VideoURL == "" {
		return domain.Video{}, errs.ErrInvalidFieldValue
	}

	oldVideo, err := s.GetVideoByID(video.ID)
	if err != nil {
		return domain.Video{}, err
	}

	if oldVideo.AuthorID != userID && userRole != domain.AdminRole {
		return domain.Video{}, errs.ErrAccessDenied
	}

	return s.repository.UpdateVideo(video)
}

func (s *Service) DeleteVideo(userID int, userRole string, videoID int) error {
	if userID <= 0 || videoID <= 0 {
		return errs.ErrInvalidFieldValue
	}

	video, err := s.GetVideoByID(videoID)
	if err != nil {
		return err
	}

	if video.AuthorID != userID && userRole != domain.AdminRole {
		return errs.ErrAccessDenied
	}

	if err = s.repository.DeleteVideo(videoID); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrVideoNotFound
		}
		return err
	}

	return nil
}
