package usecase

import (
	"errors"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func (uc *Usecase) CreateVideo(authorID int, video domain.Video) (domain.Video, error) {
	video.Title = strings.TrimSpace(video.Title)
	video.Description = strings.TrimSpace(video.Description)
	video.VideoURL = strings.TrimSpace(video.VideoURL)
	video.ThumbnailURL = strings.TrimSpace(video.ThumbnailURL)

	if authorID <= 0 || video.Title == "" || video.VideoURL == "" {
		return domain.Video{}, errs.ErrInvalidFieldValue
	}

	categoryID, err := uc.resolveCategoryIDByName(video.CategoryName)
	if err != nil {
		return domain.Video{}, err
	}

	video.CategoryID = categoryID
	video.AuthorID = authorID
	video.Status = domain.VideoStatusActive

	createdVideo, err := uc.repository.CreateVideo(video)
	if err != nil {
		return domain.Video{}, err
	}

	if video.Tags != nil {
		tagNames, err := tagNamesFromTags(video.Tags)
		if err != nil {
			return domain.Video{}, err
		}

		tags, err := uc.repository.ReplaceVideoTags(createdVideo.ID, tagNames)
		if err != nil {
			return domain.Video{}, err
		}
		createdVideo.Tags = tags
	}

	return createdVideo, nil
}

func (uc *Usecase) resolveCategoryIDByName(categoryName *string) (*int, error) {
	if categoryName == nil {
		return nil, nil
	}

	name := strings.TrimSpace(*categoryName)
	if name == "" {
		return nil, nil
	}

	category, err := uc.GetCategoryByName(name)
	if err != nil {
		return nil, err
	}

	categoryID := category.ID
	return &categoryID, nil
}

func (uc *Usecase) GetAllVideos() ([]domain.Video, error) {
	return uc.repository.GetAllVideos()
}

func (uc *Usecase) GetRecommendedVideos() ([]domain.Video, error) {
	return uc.repository.GetRecommendedVideos()
}

func (uc *Usecase) GetPersonalizedRecommendedVideos(userID int) ([]domain.Video, error) {
	if err := validateRecommendationUserID(userID); err != nil {
		return nil, err
	}
	if _, err := uc.GetUserByID(userID); err != nil {
		return nil, err
	}

	candidates, err := uc.repository.GetRecommendationCandidateVideos(recommendationCandidates)
	if err != nil {
		return nil, err
	}

	profile, err := uc.buildRecommendationProfile(userID, candidates)
	if err != nil {
		return nil, err
	}

	scoredCandidates := buildRecommendationCandidates(candidates, profile)
	return uc.recommendationEngine.Recommend(profile, scoredCandidates, recommendationLimit), nil
}

func (uc *Usecase) GetVideoByID(id int) (domain.Video, error) {
	if id <= 0 {
		return domain.Video{}, errs.ErrInvalidFieldValue
	}

	video, err := uc.repository.GetVideoByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.Video{}, errs.ErrVideoNotFound
		}
		return domain.Video{}, err
	}
	qualities, err := uc.repository.GetVideoQualities(id)
	if err != nil {
		return domain.Video{}, err
	}
	video.Qualities = qualities
	tags, err := uc.repository.GetVideoTags(id)
	if err != nil {
		return domain.Video{}, err
	}
	video.Tags = tags

	return video, nil
}

func (uc *Usecase) GetUserVideos(userID int) ([]domain.Video, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetUserByID(userID); err != nil {
		return nil, err
	}

	return uc.repository.GetVideosByAuthorID(userID)
}

func (uc *Usecase) IncrementVideoViews(id int) error {
	if id <= 0 {
		return errs.ErrInvalidFieldValue
	}

	if err := uc.repository.IncrementVideoViews(id); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrVideoNotFound
		}
		return err
	}

	return nil
}

func (uc *Usecase) UpdateVideo(userID int, userRole string, video domain.Video) (domain.Video, error) {
	video.Title = strings.TrimSpace(video.Title)
	video.Description = strings.TrimSpace(video.Description)
	video.ThumbnailURL = strings.TrimSpace(video.ThumbnailURL)

	if userID <= 0 || video.ID <= 0 || video.Title == "" {
		return domain.Video{}, errs.ErrInvalidFieldValue
	}

	oldVideo, err := uc.GetVideoByID(video.ID)
	if err != nil {
		return domain.Video{}, err
	}

	if oldVideo.AuthorID != userID && userRole != domain.AdminRole {
		return domain.Video{}, errs.ErrAccessDenied
	}

	categoryID := oldVideo.CategoryID
	if video.CategoryName != nil {
		categoryID, err = uc.resolveCategoryIDByName(video.CategoryName)
		if err != nil {
			return domain.Video{}, err
		}
	}

	video.VideoURL = oldVideo.VideoURL
	video.CategoryID = categoryID
	if video.ThumbnailURL == "" {
		video.ThumbnailURL = oldVideo.ThumbnailURL
	}

	updatedVideo, err := uc.repository.UpdateVideo(video)
	if err != nil {
		return domain.Video{}, err
	}

	if video.Tags != nil {
		tagNames, err := tagNamesFromTags(video.Tags)
		if err != nil {
			return domain.Video{}, err
		}

		tags, err := uc.repository.ReplaceVideoTags(updatedVideo.ID, tagNames)
		if err != nil {
			return domain.Video{}, err
		}
		updatedVideo.Tags = tags
	} else {
		updatedVideo.Tags = oldVideo.Tags
	}

	return updatedVideo, nil
}

func (uc *Usecase) DeleteVideo(userID int, userRole string, videoID int) error {
	if userID <= 0 || videoID <= 0 {
		return errs.ErrInvalidFieldValue
	}

	video, err := uc.GetVideoByID(videoID)
	if err != nil {
		return err
	}

	if video.AuthorID != userID && userRole != domain.AdminRole {
		return errs.ErrAccessDenied
	}

	if err = uc.repository.DeleteVideo(videoID); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrVideoNotFound
		}
		return err
	}

	return nil
}

func (uc *Usecase) CategoryExists(categoryName *string) error {
	if categoryName == nil {
		return nil
	}

	if *categoryName == "" {
		return errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetCategoryByName(*categoryName); err != nil {
		return err
	}

	return nil
}

func (uc *Usecase) SearchVideosByTitle(title string) ([]domain.Video, error) {
	return uc.SearchVideosByTitleForUser(nil, title)
}

func (uc *Usecase) SearchVideosByTitleForUser(userID *int, title string) ([]domain.Video, error) {
	title = strings.TrimSpace(title)
	if title == "" {
		return uc.repository.GetRecommendedVideos()
	}

	if userID != nil {
		if *userID <= 0 {
			return nil, errs.ErrInvalidFieldValue
		}
		if err := uc.repository.SaveVideoSearchHistory(domain.VideoSearchHistory{
			UserID: *userID,
			Query:  title,
		}); err != nil {
			return nil, err
		}
	}

	return uc.repository.SearchVideosByTitle(title)
}

func (uc *Usecase) SaveVideoWatchProgress(userID int, videoID int, watchedSeconds int, durationSeconds int) (domain.VideoWatchHistory, error) {
	if userID <= 0 || videoID <= 0 || watchedSeconds < 0 || durationSeconds <= 0 {
		return domain.VideoWatchHistory{}, errs.ErrInvalidFieldValue
	}

	if _, err := uc.repository.GetVideoByID(videoID); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.VideoWatchHistory{}, errs.ErrVideoNotFound
		}
		return domain.VideoWatchHistory{}, err
	}

	watchedPercent := calculateWatchedPercent(watchedSeconds, durationSeconds)
	return uc.repository.SaveVideoWatchProgress(domain.VideoWatchHistory{
		UserID:          userID,
		VideoID:         videoID,
		WatchedSeconds:  watchedSeconds,
		DurationSeconds: durationSeconds,
		WatchedPercent:  watchedPercent,
		IsCompleted:     watchedPercent > 80,
	})
}

func (uc *Usecase) GenerateVideoQualities(userID int, userRole string, videoID int, inputPath string, outputDir string, outputURLPrefix string) ([]domain.VideoQuality, error) {
	if userID <= 0 || videoID <= 0 || inputPath == "" || outputDir == "" || outputURLPrefix == "" {
		return nil, errs.ErrInvalidFieldValue
	}

	video, err := uc.GetVideoByID(videoID)
	if err != nil {
		return nil, err
	}

	if video.AuthorID != userID && userRole != domain.AdminRole {
		return nil, errs.ErrAccessDenied
	}

	qualities, err := uc.ffmpegSettings.GenerateVideoQualities(inputPath, outputDir, outputURLPrefix)
	if err != nil {
		return nil, err
	}

	savedQualities := make([]domain.VideoQuality, 0, len(qualities))

	for _, quality := range qualities {
		quality.VideoID = videoID

		savedQuality, err := uc.repository.CreateVideoQuality(quality)
		if err != nil {
			return nil, err
		}

		savedQualities = append(savedQualities, savedQuality)
	}

	return savedQualities, nil
}

func (uc *Usecase) GetPlaybackSpeeds() []float64 {
	return []float64{0.25, 1.0, 1.25, 1.5, 2.0}
}

func (uc *Usecase) DeleteVideoWithArchive(userID int, userRole string, videoID int, inputPath string, archivePath string, archiveURL string) error {
	if userID <= 0 || videoID <= 0 || inputPath == "" || archivePath == "" || archiveURL == "" {
		return errs.ErrInvalidFieldValue
	}

	video, err := uc.GetVideoByID(videoID)
	if err != nil {
		return err
	}

	if video.AuthorID != userID && userRole != domain.AdminRole {
		return errs.ErrAccessDenied
	}

	if err = uc.ffmpegSettings.GenerateArchive144p(inputPath, archivePath); err != nil {
		return err
	}

	if err = uc.repository.ArchiveDeletedVideo(videoID, archiveURL); err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return errs.ErrVideoNotFound
		}
		return err
	}

	return nil
}
