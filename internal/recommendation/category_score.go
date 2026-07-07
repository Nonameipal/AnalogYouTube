package recommendation

import "github.com/Nonameipal/AnalogYouTube/internal/domain"

type CategoryScorer struct{}

func (CategoryScorer) Name() string {
	return "category_score"
}

func (CategoryScorer) Score(ctx Context, candidate domain.RecommendationCandidate) float64 {
	if candidate.Video.CategoryID == nil {
		return 0
	}

	score := normalized(ctx.CategoryAffinity[*candidate.Video.CategoryID], ctx.MaxCategoryAffinity)
	return score * 12
}
