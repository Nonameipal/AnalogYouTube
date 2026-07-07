package recommendation

import "github.com/Nonameipal/AnalogYouTube/internal/domain"

type SubscriptionScorer struct{}

func (SubscriptionScorer) Name() string {
	return "subscription_score"
}

func (SubscriptionScorer) Score(ctx Context, candidate domain.RecommendationCandidate) float64 {
	if _, ok := ctx.SubscribedAuthorIDs[candidate.Video.AuthorID]; ok {
		return 16
	}

	return 0
}
