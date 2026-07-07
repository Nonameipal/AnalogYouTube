package usecase

import (
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

const (
	recommendationLimit      = 20
	recommendationCandidates = 200
	recommendationHistory    = 200
	recommendationSearches   = 100
)

func (uc *Usecase) buildRecommendationProfile(userID int, candidates []domain.Video) (domain.RecommendationProfile, error) {
	watchHistory, err := uc.repository.GetUserWatchHistory(userID, recommendationHistory)
	if err != nil {
		return domain.RecommendationProfile{}, err
	}

	searchHistory, err := uc.repository.GetUserSearchHistory(userID, recommendationSearches)
	if err != nil {
		return domain.RecommendationProfile{}, err
	}

	likedVideoIDs, err := uc.repository.GetUserLikedVideoIDs(userID)
	if err != nil {
		return domain.RecommendationProfile{}, err
	}

	commentedVideoIDs, err := uc.repository.GetUserCommentedVideoIDs(userID)
	if err != nil {
		return domain.RecommendationProfile{}, err
	}

	playlistVideoIDs, err := uc.repository.GetUserPlaylistVideoIDs(userID)
	if err != nil {
		return domain.RecommendationProfile{}, err
	}

	subscribedAuthorIDs, err := uc.repository.GetSubscribedAuthorIDs(userID)
	if err != nil {
		return domain.RecommendationProfile{}, err
	}

	videoIDs := collectRecommendationVideoIDs(candidates, watchHistory, likedVideoIDs, commentedVideoIDs, playlistVideoIDs)
	videosByID, err := uc.getVideosByID(videoIDs)
	if err != nil {
		return domain.RecommendationProfile{}, err
	}

	tagsByVideoID, err := uc.repository.GetTagsByVideoIDs(videoIDs)
	if err != nil {
		return domain.RecommendationProfile{}, err
	}

	candidateIDs := videoIDsFromVideos(candidates)
	popularities, err := uc.repository.GetVideoPopularity(candidateIDs)
	if err != nil {
		return domain.RecommendationProfile{}, err
	}

	return domain.RecommendationProfile{
		UserID:                userID,
		WatchHistory:          watchHistory,
		SearchHistory:         searchHistory,
		LikedVideoIDs:         likedVideoIDs,
		CommentedVideoIDs:     commentedVideoIDs,
		PlaylistVideoIDs:      playlistVideoIDs,
		SubscribedAuthorIDs:   subscribedAuthorIDs,
		VideosByID:            videosByID,
		TagsByVideoID:         tagsByVideoID,
		CandidatePopularities: popularities,
	}, nil
}

func (uc *Usecase) getVideosByID(videoIDs []int) (map[int]domain.Video, error) {
	videos, err := uc.repository.GetVideosByIDs(videoIDs)
	if err != nil {
		return nil, err
	}

	videosByID := make(map[int]domain.Video, len(videos))
	for _, video := range videos {
		videosByID[video.ID] = video
	}

	return videosByID, nil
}

func collectRecommendationVideoIDs(candidates []domain.Video, watchHistory []domain.VideoWatchHistory, likedIDs []int, commentedIDs []int, playlistIDs []int) []int {
	seen := make(map[int]struct{})
	ids := make([]int, 0, len(candidates)+len(watchHistory)+len(likedIDs)+len(commentedIDs)+len(playlistIDs))

	add := func(id int) {
		if id <= 0 {
			return
		}
		if _, ok := seen[id]; ok {
			return
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}

	for _, video := range candidates {
		add(video.ID)
	}
	for _, item := range watchHistory {
		add(item.VideoID)
	}
	for _, id := range likedIDs {
		add(id)
	}
	for _, id := range commentedIDs {
		add(id)
	}
	for _, id := range playlistIDs {
		add(id)
	}

	return ids
}

func videoIDsFromVideos(videos []domain.Video) []int {
	ids := make([]int, 0, len(videos))
	for _, video := range videos {
		ids = append(ids, video.ID)
	}

	return ids
}

func buildRecommendationCandidates(videos []domain.Video, profile domain.RecommendationProfile) []domain.RecommendationCandidate {
	candidates := make([]domain.RecommendationCandidate, 0, len(videos))
	for _, video := range videos {
		candidates = append(candidates, domain.RecommendationCandidate{
			Video:      video,
			Tags:       profile.TagsByVideoID[video.ID],
			Popularity: profile.CandidatePopularities[video.ID],
		})
	}

	return candidates
}

func validateRecommendationUserID(userID int) error {
	if userID <= 0 {
		return errs.ErrInvalidFieldValue
	}

	return nil
}
