package recommendation

import (
	"sort"
	"time"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
)

const defaultRecommendationLimit = 20

type Scorer interface {
	Name() string
	Score(ctx Context, candidate domain.RecommendationCandidate) float64
}

type Engine struct {
	scorers []Scorer
	now     func() time.Time
}

type scoredCandidate struct {
	Candidate domain.RecommendationCandidate
	Score     float64
}

func NewEngine() *Engine {
	return &Engine{
		scorers: []Scorer{
			CategoryScorer{},
			TagScorer{},
			WatchScorer{},
			SearchScorer{},
			PlaylistScorer{},
			SubscriptionScorer{},
			CommentScorer{},
			LikeScorer{},
			PopularityScorer{},
			FreshnessScorer{},
		},
		now: time.Now,
	}
}

func (e *Engine) Recommend(profile domain.RecommendationProfile, candidates []domain.RecommendationCandidate, limit int) []domain.Video {
	if limit <= 0 {
		limit = defaultRecommendationLimit
	}

	ctx := newContext(profile, candidates, e.now())
	scored := make([]scoredCandidate, 0, len(candidates))
	for _, candidate := range candidates {
		scored = append(scored, scoredCandidate{
			Candidate: candidate,
			Score:     e.scoreCandidate(ctx, candidate),
		})
	}

	sort.SliceStable(scored, func(i, j int) bool {
		if scored[i].Score == scored[j].Score {
			return scored[i].Candidate.Video.CreatedAt.After(scored[j].Candidate.Video.CreatedAt)
		}
		return scored[i].Score > scored[j].Score
	})

	diversified := applyDiversity(scored, limit)
	videos := make([]domain.Video, 0, len(diversified))
	for _, item := range diversified {
		video := item.Candidate.Video
		video.Tags = item.Candidate.Tags
		videos = append(videos, video)
	}

	return videos
}

func (e *Engine) scoreCandidate(ctx Context, candidate domain.RecommendationCandidate) float64 {
	var total float64
	for _, scorer := range e.scorers {
		total += scorer.Score(ctx, candidate)
	}

	return total
}
