package recommendation

import "github.com/Nonameipal/AnalogYouTube/internal/domain"

type CommentScorer struct{}

func (CommentScorer) Name() string {
	return "comment_score"
}

func (CommentScorer) Score(ctx Context, candidate domain.RecommendationCandidate) float64 {
	var score float64
	if _, ok := ctx.CommentedVideoIDs[candidate.Video.ID]; ok {
		score += 5
	}

	if candidate.Video.CategoryID != nil {
		score += normalized(ctx.CommentCategoryAffinity[*candidate.Video.CategoryID], ctx.MaxCommentCategory) * 8
	}

	if len(candidate.Tags) > 0 {
		var tagScore float64
		for _, tag := range candidate.Tags {
			tagScore += normalized(ctx.CommentTagAffinity[normalizeText(tag.Name)], ctx.MaxCommentTag)
		}
		score += tagScore / float64(len(candidate.Tags)) * 10
	}

	return score
}
