package recommendation

import "github.com/Nonameipal/AnalogYouTube/internal/domain"

type TagScorer struct{}

func (TagScorer) Name() string {
	return "tag_score"
}

func (TagScorer) Score(ctx Context, candidate domain.RecommendationCandidate) float64 {
	if len(candidate.Tags) == 0 {
		return 0
	}

	var score float64
	for _, tag := range candidate.Tags {
		score += normalized(ctx.TagAffinity[normalizeText(tag.Name)], ctx.MaxTagAffinity)
	}

	return score / float64(len(candidate.Tags)) * 18
}
