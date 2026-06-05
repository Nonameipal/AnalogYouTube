package service

import "github.com/Nonameipal/AnalogYouTube/internal/errs"

func (s *Service) LikeVideo(userID int, videoID int) error {
	if userID <= 0 || videoID <= 0 {
		return errs.ErrInvalidFieldValue
	}

	if _, err := s.GetVideoByID(videoID); err != nil {
		return err
	}

	return s.repository.LikeVideo(userID, videoID)
}

func (s *Service) UnlikeVideo(userID int, videoID int) error {
	if userID <= 0 || videoID <= 0 {
		return errs.ErrInvalidFieldValue
	}

	if _, err := s.GetVideoByID(videoID); err != nil {
		return err
	}

	return s.repository.UnlikeVideo(userID, videoID)
}

func (s *Service) GetVideoLikesCount(videoID int) (int, error) {
	if videoID <= 0 {
		return 0, errs.ErrInvalidFieldValue
	}

	if _, err := s.GetVideoByID(videoID); err != nil {
		return 0, err
	}

	return s.repository.GetVideoLikesCount(videoID)
}

func (s *Service) IsVideoLikedByUser(userID int, videoID int) (bool, error) {
	if userID <= 0 || videoID <= 0 {
		return false, errs.ErrInvalidFieldValue
	}

	if _, err := s.GetVideoByID(videoID); err != nil {
		return false, err
	}

	return s.repository.IsVideoLikedByUser(userID, videoID)
}
