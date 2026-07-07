package recommendation

import (
	"math"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
)

type PopularityScorer struct{}

func (PopularityScorer) Name() string {
	return "popularity_score"
}

func (PopularityScorer) Score(ctx Context, candidate domain.RecommendationCandidate) float64 {
	viewsScore := normalizedLogFloat(float64(candidate.Video.Views), float64(ctx.MaxViews)) * 5
	likesScore := normalizedLogFloat(float64(candidate.Popularity.LikesCount), float64(ctx.MaxLikes)) * 5
	commentsScore := normalizedLogFloat(float64(candidate.Popularity.CommentsCount), float64(ctx.MaxComments)) * 6
	playlistScore := normalizedLogFloat(float64(candidate.Popularity.PlaylistAddsCount), float64(ctx.MaxPlaylistAdds)) * 4

	return viewsScore + likesScore + commentsScore + playlistScore
}

func normalizedLogFloat(value float64, max float64) float64 {
	if value <= 0 || max <= 0 {
		return 0
	}

	return math.Log1p(value) / math.Log1p(max)
}
