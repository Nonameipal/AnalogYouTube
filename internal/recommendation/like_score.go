package recommendation

import "github.com/Nonameipal/AnalogYouTube/internal/domain"

type LikeScorer struct{}

func (LikeScorer) Name() string {
	return "like_score"
}

func (LikeScorer) Score(ctx Context, candidate domain.RecommendationCandidate) float64 {
	var score float64
	if _, ok := ctx.LikedVideoIDs[candidate.Video.ID]; ok {
		score += 3
	}

	if candidate.Video.CategoryID != nil {
		score += normalized(ctx.LikeCategoryAffinity[*candidate.Video.CategoryID], ctx.MaxLikeCategory) * 6
	}

	if len(candidate.Tags) > 0 {
		var tagScore float64
		for _, tag := range candidate.Tags {
			tagScore += normalized(ctx.LikeTagAffinity[normalizeText(tag.Name)], ctx.MaxLikeTag)
		}
		score += tagScore / float64(len(candidate.Tags)) * 8
	}

	return score
}
