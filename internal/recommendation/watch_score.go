package recommendation

import "github.com/Nonameipal/AnalogYouTube/internal/domain"

type WatchScorer struct{}

func (WatchScorer) Name() string {
	return "watch_score"
}

func (WatchScorer) Score(ctx Context, candidate domain.RecommendationCandidate) float64 {
	score := watchedVideoPenalty(ctx, candidate.Video.ID)
	if candidate.Video.CategoryID != nil {
		score += normalized(ctx.WatchCategoryAffinity[*candidate.Video.CategoryID], ctx.MaxWatchCategory) * 8
	}

	if len(candidate.Tags) > 0 {
		var tagScore float64
		for _, tag := range candidate.Tags {
			tagScore += normalized(ctx.WatchTagAffinity[normalizeText(tag.Name)], ctx.MaxWatchTag)
		}
		score += tagScore / float64(len(candidate.Tags)) * 10
	}

	return score
}

func watchedVideoPenalty(ctx Context, videoID int) float64 {
	watchedPercent := ctx.WatchedPercentByVideoID[videoID]
	switch {
	case watchedPercent >= 90:
		return -1000
	case watchedPercent >= 80:
		return -40
	case watchedPercent > 0 && watchedPercent < 20:
		return -4
	default:
		return 0
	}
}
