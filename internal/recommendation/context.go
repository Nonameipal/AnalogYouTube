package recommendation

import (
	"strings"
	"time"
	"unicode"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
)

type Context struct {
	Profile domain.RecommendationProfile
	Now     time.Time

	WatchedPercentByVideoID map[int]float64
	LikedVideoIDs           map[int]struct{}
	CommentedVideoIDs       map[int]struct{}
	PlaylistVideoIDs        map[int]struct{}
	SubscribedAuthorIDs     map[int]struct{}

	CategoryAffinity map[int]float64
	TagAffinity      map[string]float64

	WatchCategoryAffinity    map[int]float64
	WatchTagAffinity         map[string]float64
	LikeCategoryAffinity     map[int]float64
	LikeTagAffinity          map[string]float64
	CommentCategoryAffinity  map[int]float64
	CommentTagAffinity       map[string]float64
	PlaylistCategoryAffinity map[int]float64
	PlaylistTagAffinity      map[string]float64

	SearchTerms map[string]int

	MaxCategoryAffinity float64
	MaxTagAffinity      float64
	MaxWatchCategory    float64
	MaxWatchTag         float64
	MaxLikeCategory     float64
	MaxLikeTag          float64
	MaxCommentCategory  float64
	MaxCommentTag       float64
	MaxPlaylistCategory float64
	MaxPlaylistTag      float64

	MaxViews        int64
	MaxLikes        int
	MaxComments     int
	MaxPlaylistAdds int
	HasUserSignals  bool
}

func newContext(profile domain.RecommendationProfile, candidates []domain.RecommendationCandidate, now time.Time) Context {
	ctx := Context{
		Profile:                  profile,
		Now:                      now,
		WatchedPercentByVideoID:  make(map[int]float64),
		LikedVideoIDs:            intSet(profile.LikedVideoIDs),
		CommentedVideoIDs:        intSet(profile.CommentedVideoIDs),
		PlaylistVideoIDs:         intSet(profile.PlaylistVideoIDs),
		SubscribedAuthorIDs:      intSet(profile.SubscribedAuthorIDs),
		CategoryAffinity:         make(map[int]float64),
		TagAffinity:              make(map[string]float64),
		WatchCategoryAffinity:    make(map[int]float64),
		WatchTagAffinity:         make(map[string]float64),
		LikeCategoryAffinity:     make(map[int]float64),
		LikeTagAffinity:          make(map[string]float64),
		CommentCategoryAffinity:  make(map[int]float64),
		CommentTagAffinity:       make(map[string]float64),
		PlaylistCategoryAffinity: make(map[int]float64),
		PlaylistTagAffinity:      make(map[string]float64),
		SearchTerms:              make(map[string]int),
	}

	ctx.HasUserSignals = hasUserSignals(profile)
	ctx.addWatchSignals()
	ctx.addVideoIDSignals(profile.LikedVideoIDs, ctx.LikeCategoryAffinity, ctx.LikeTagAffinity, 2.0)
	ctx.addVideoIDSignals(profile.CommentedVideoIDs, ctx.CommentCategoryAffinity, ctx.CommentTagAffinity, 3.0)
	ctx.addVideoIDSignals(profile.PlaylistVideoIDs, ctx.PlaylistCategoryAffinity, ctx.PlaylistTagAffinity, 2.5)
	ctx.addSearchSignals()
	ctx.calculateMaxAffinities()
	ctx.calculateMaxPopularity(candidates)

	return ctx
}

func (ctx *Context) addWatchSignals() {
	for _, item := range ctx.Profile.WatchHistory {
		weight := clamp(item.WatchedPercent/100, 0.1, 1.0)
		if item.IsCompleted {
			weight += 0.4
		}

		if current := ctx.WatchedPercentByVideoID[item.VideoID]; item.WatchedPercent > current {
			ctx.WatchedPercentByVideoID[item.VideoID] = item.WatchedPercent
		}

		ctx.addVideoSignal(item.VideoID, ctx.WatchCategoryAffinity, ctx.WatchTagAffinity, weight)
	}
}

func (ctx *Context) addVideoIDSignals(videoIDs []int, categoryAffinity map[int]float64, tagAffinity map[string]float64, weight float64) {
	for _, videoID := range videoIDs {
		ctx.addVideoSignal(videoID, categoryAffinity, tagAffinity, weight)
	}
}

func (ctx *Context) addVideoSignal(videoID int, categoryAffinity map[int]float64, tagAffinity map[string]float64, weight float64) {
	video, ok := ctx.Profile.VideosByID[videoID]
	if !ok {
		return
	}

	if video.CategoryID != nil {
		categoryAffinity[*video.CategoryID] += weight
		ctx.CategoryAffinity[*video.CategoryID] += weight
	}

	for _, tag := range ctx.Profile.TagsByVideoID[videoID] {
		key := normalizeText(tag.Name)
		if key == "" {
			continue
		}
		tagAffinity[key] += weight
		ctx.TagAffinity[key] += weight
	}
}

func (ctx *Context) addSearchSignals() {
	for _, history := range ctx.Profile.SearchHistory {
		for _, term := range searchTerms(history.Query) {
			ctx.SearchTerms[term]++
		}
	}
}

func (ctx *Context) calculateMaxAffinities() {
	ctx.MaxCategoryAffinity = maxFloatMapValue(ctx.CategoryAffinity)
	ctx.MaxTagAffinity = maxFloatMapValueString(ctx.TagAffinity)
	ctx.MaxWatchCategory = maxFloatMapValue(ctx.WatchCategoryAffinity)
	ctx.MaxWatchTag = maxFloatMapValueString(ctx.WatchTagAffinity)
	ctx.MaxLikeCategory = maxFloatMapValue(ctx.LikeCategoryAffinity)
	ctx.MaxLikeTag = maxFloatMapValueString(ctx.LikeTagAffinity)
	ctx.MaxCommentCategory = maxFloatMapValue(ctx.CommentCategoryAffinity)
	ctx.MaxCommentTag = maxFloatMapValueString(ctx.CommentTagAffinity)
	ctx.MaxPlaylistCategory = maxFloatMapValue(ctx.PlaylistCategoryAffinity)
	ctx.MaxPlaylistTag = maxFloatMapValueString(ctx.PlaylistTagAffinity)
}

func (ctx *Context) calculateMaxPopularity(candidates []domain.RecommendationCandidate) {
	for _, candidate := range candidates {
		if candidate.Video.Views > ctx.MaxViews {
			ctx.MaxViews = candidate.Video.Views
		}
		if candidate.Popularity.LikesCount > ctx.MaxLikes {
			ctx.MaxLikes = candidate.Popularity.LikesCount
		}
		if candidate.Popularity.CommentsCount > ctx.MaxComments {
			ctx.MaxComments = candidate.Popularity.CommentsCount
		}
		if candidate.Popularity.PlaylistAddsCount > ctx.MaxPlaylistAdds {
			ctx.MaxPlaylistAdds = candidate.Popularity.PlaylistAddsCount
		}
	}
}

func hasUserSignals(profile domain.RecommendationProfile) bool {
	return len(profile.WatchHistory) > 0 ||
		len(profile.SearchHistory) > 0 ||
		len(profile.LikedVideoIDs) > 0 ||
		len(profile.CommentedVideoIDs) > 0 ||
		len(profile.PlaylistVideoIDs) > 0 ||
		len(profile.SubscribedAuthorIDs) > 0
}

func intSet(ids []int) map[int]struct{} {
	set := make(map[int]struct{}, len(ids))
	for _, id := range ids {
		set[id] = struct{}{}
	}

	return set
}

func searchTerms(query string) []string {
	terms := make([]string, 0)
	query = normalizeText(query)
	if query == "" {
		return terms
	}

	terms = append(terms, query)
	for _, token := range strings.FieldsFunc(query, func(r rune) bool {
		return unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r)
	}) {
		if token = normalizeText(token); token != "" {
			terms = append(terms, token)
			if token == "golang" {
				terms = append(terms, "go")
			}
		}
	}

	return terms
}

func normalizeText(value string) string {
	return strings.ToLower(strings.TrimSpace(value))
}

func maxFloatMapValue(values map[int]float64) float64 {
	var max float64
	for _, value := range values {
		if value > max {
			max = value
		}
	}

	return max
}

func maxFloatMapValueString(values map[string]float64) float64 {
	var max float64
	for _, value := range values {
		if value > max {
			max = value
		}
	}

	return max
}

func clamp(value float64, min float64, max float64) float64 {
	if value < min {
		return min
	}
	if value > max {
		return max
	}

	return value
}

func normalized(value float64, max float64) float64 {
	if max <= 0 || value <= 0 {
		return 0
	}

	return value / max
}
