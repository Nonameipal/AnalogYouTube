package recommendation

import "github.com/Nonameipal/AnalogYouTube/internal/domain"

type PlaylistScorer struct{}

func (PlaylistScorer) Name() string {
	return "playlist_score"
}

func (PlaylistScorer) Score(ctx Context, candidate domain.RecommendationCandidate) float64 {
	var score float64
	if _, ok := ctx.PlaylistVideoIDs[candidate.Video.ID]; ok {
		score += 4
	}

	if candidate.Video.CategoryID != nil {
		score += normalized(ctx.PlaylistCategoryAffinity[*candidate.Video.CategoryID], ctx.MaxPlaylistCategory) * 7
	}

	if len(candidate.Tags) > 0 {
		var tagScore float64
		for _, tag := range candidate.Tags {
			tagScore += normalized(ctx.PlaylistTagAffinity[normalizeText(tag.Name)], ctx.MaxPlaylistTag)
		}
		score += tagScore / float64(len(candidate.Tags)) * 8
	}

	return score
}
