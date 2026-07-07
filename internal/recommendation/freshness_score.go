package recommendation

import "github.com/Nonameipal/AnalogYouTube/internal/domain"

type FreshnessScorer struct{}

func (FreshnessScorer) Name() string {
	return "freshness_score"
}

func (FreshnessScorer) Score(ctx Context, candidate domain.RecommendationCandidate) float64 {
	ageDays := ctx.Now.Sub(candidate.Video.CreatedAt).Hours() / 24
	switch {
	case ageDays < 0:
		return 7
	case ageDays <= 3:
		return 7
	case ageDays <= 7:
		return 5
	case ageDays <= 30:
		return 3
	case ageDays <= 90:
		return 1
	default:
		return 0
	}
}
