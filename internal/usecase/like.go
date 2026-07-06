package usecase

import "github.com/Nonameipal/AnalogYouTube/internal/errs"

func (uc *Usecase) LikeVideo(userID int, videoID int) error {
	if userID <= 0 || videoID <= 0 {
		return errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetVideoByID(videoID); err != nil {
		return err
	}

	return uc.repository.LikeVideo(userID, videoID)
}

func (uc *Usecase) UnlikeVideo(userID int, videoID int) error {
	if userID <= 0 || videoID <= 0 {
		return errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetVideoByID(videoID); err != nil {
		return err
	}

	return uc.repository.UnlikeVideo(userID, videoID)
}

func (uc *Usecase) GetVideoLikesCount(videoID int) (int, error) {
	if videoID <= 0 {
		return 0, errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetVideoByID(videoID); err != nil {
		return 0, err
	}

	return uc.repository.GetVideoLikesCount(videoID)
}

func (uc *Usecase) IsVideoLikedByUser(userID int, videoID int) (bool, error) {
	if userID <= 0 || videoID <= 0 {
		return false, errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetVideoByID(videoID); err != nil {
		return false, err
	}

	return uc.repository.IsVideoLikedByUser(userID, videoID)
}
