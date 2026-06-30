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

	categoryID, err := s.resolveCategoryIDByName(video.CategoryName)
    if err != nil {
	return domain.Video{}, err
}


	video.CategoryID = categoryID
	video.AuthorID = authorID
	video.Status = domain.VideoStatusActive

	return s.repository.CreateVideo(video)
}

func (s *Service) resolveCategoryIDByName(categoryName *string) (*int, error) {
	if categoryName == nil {
		return nil, nil
	}

	name := strings.TrimSpace(*categoryName)
	if name == "" {
		return nil, nil
	}

	category, err := s.GetCategoryByName(name)
	if err != nil {
		return nil, err
	}

	categoryID := category.ID
	return &categoryID, nil
}

func (s *Service) GetAllVideos() ([]domain.Video, error) {
	return s.repository.GetAllVideos()
}

func (s *Service) GetRecommendedVideos() ([]domain.Video, error) {
	return s.repository.GetRecommendedVideos()
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

func (s *Service) GetUserVideos(userID int) ([]domain.Video, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	if _, err := s.GetUserByID(userID); err != nil {
		return nil, err
	}

	return s.repository.GetVideosByAuthorID(userID)
}

func (s *Service) IncrementVideoViews(id int) error {
	if id <= 0 {
		return errs.ErrInvalidFieldValue
	}

	if err := s.repository.IncrementVideoViews(id); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrVideoNotFound
		}
		return err
	}

	return nil
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

	categoryID := oldVideo.CategoryID
    if video.CategoryName != nil {
	categoryID, err = s.resolveCategoryIDByName(video.CategoryName)
	if err != nil {
		return domain.Video{}, err
	}
}

video.CategoryID = categoryID

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

func (s *Service) CategoryExists(categoryName *string) error {
	if categoryName == nil {
		return nil
	}

	if *categoryName == "" {
		return errs.ErrInvalidFieldValue
	}

	if _, err := s.GetCategoryByName(*categoryName); err != nil {
		return err
	}

	return nil
}

func (s *Service) SearchVideosByTitle(title string) ([]domain.Video, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return s.repository.GetRecommendedVideos()
	}

	return s.repository.SearchVideosByTitle(title)
}
