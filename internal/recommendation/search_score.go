package recommendation

import (
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
)

type SearchScorer struct{}

func (SearchScorer) Name() string {
	return "search_score"
}

func (SearchScorer) Score(ctx Context, candidate domain.RecommendationCandidate) float64 {
	if len(ctx.SearchTerms) == 0 {
		return 0
	}

	text := searchableCandidateText(candidate)
	var score float64
	for term, count := range ctx.SearchTerms {
		if term == "" || !strings.Contains(text, term) {
			continue
		}
		score += clamp(float64(count), 1, 5) * 2.5
	}

	return clamp(score, 0, 18)
}

func searchableCandidateText(candidate domain.RecommendationCandidate) string {
	parts := []string{
		candidate.Video.Title,
		candidate.Video.Description,
	}
	for _, tag := range candidate.Tags {
		parts = append(parts, tag.Name)
	}

	return normalizeText(strings.Join(parts, " "))
}
