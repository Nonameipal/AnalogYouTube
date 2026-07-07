package domain

type VideoPopularity struct {
	VideoID           int `json:"video_id"`
	LikesCount        int `json:"likes_count"`
	CommentsCount     int `json:"comments_count"`
	PlaylistAddsCount int `json:"playlist_adds_count"`
}

type RecommendationProfile struct {
	UserID                int
	WatchHistory          []VideoWatchHistory
	SearchHistory         []VideoSearchHistory
	LikedVideoIDs         []int
	CommentedVideoIDs     []int
	PlaylistVideoIDs      []int
	SubscribedAuthorIDs   []int
	VideosByID            map[int]Video
	TagsByVideoID         map[int][]Tag
	CandidatePopularities map[int]VideoPopularity
}

type RecommendationCandidate struct {
	Video      Video
	Tags       []Tag
	Popularity VideoPopularity
}
