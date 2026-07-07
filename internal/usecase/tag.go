package usecase

import (
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

const (
	maxVideoTags     = 20
	maxTagNameLength = 100
)

func (uc *Usecase) GetAllTags() ([]domain.Tag, error) {
	return uc.repository.GetAllTags()
}

func (uc *Usecase) GetVideoTags(videoID int) ([]domain.Tag, error) {
	if videoID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	if _, err := uc.repository.GetVideoByID(videoID); err != nil {
		return nil, err
	}

	return uc.repository.GetVideoTags(videoID)
}

func (uc *Usecase) UpdateVideoTags(userID int, userRole string, videoID int, tagNames []string) ([]domain.Tag, error) {
	if userID <= 0 || videoID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	video, err := uc.GetVideoByID(videoID)
	if err != nil {
		return nil, err
	}

	if video.AuthorID != userID && userRole != domain.AdminRole {
		return nil, errs.ErrAccessDenied
	}

	normalizedNames, err := normalizeTagNames(tagNames)
	if err != nil {
		return nil, err
	}

	return uc.repository.ReplaceVideoTags(videoID, normalizedNames)
}

func tagNamesFromTags(tags []domain.Tag) ([]string, error) {
	names := make([]string, 0, len(tags))
	for _, tag := range tags {
		names = append(names, tag.Name)
	}

	return normalizeTagNames(names)
}

func normalizeTagNames(tagNames []string) ([]string, error) {
	if len(tagNames) > maxVideoTags {
		return nil, errs.ErrInvalidFieldValue
	}

	seen := make(map[string]struct{}, len(tagNames))
	normalized := make([]string, 0, len(tagNames))
	for _, name := range tagNames {
		name = normalizeTagName(name)
		if name == "" {
			continue
		}
		if len(name) > maxTagNameLength {
			return nil, errs.ErrInvalidFieldValue
		}

		key := strings.ToLower(name)
		if _, ok := seen[key]; ok {
			continue
		}

		seen[key] = struct{}{}
		normalized = append(normalized, name)
	}

	return normalized, nil
}

func normalizeTagName(name string) string {
	fields := strings.Fields(strings.TrimSpace(name))
	return strings.Join(fields, " ")
}

func calculateWatchedPercent(watchedSeconds int, durationSeconds int) float64 {
	percent := float64(watchedSeconds) / float64(durationSeconds) * 100
	if percent > 100 {
		return 100
	}

	return percent
}
