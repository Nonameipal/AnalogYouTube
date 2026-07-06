package ports

import "github.com/Nonameipal/AnalogYouTube/internal/domain"

type RepositoryI interface {
	CreateUser(user domain.User) error
	GetUserByUsername(username string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	GetUserByID(id int) (domain.User, error)
	UpdateUserProfile(user domain.User) (domain.User, error)

	CreateVideo(video domain.Video) (domain.Video, error)
	GetAllVideos() ([]domain.Video, error)
	GetRecommendedVideos() ([]domain.Video, error)
	GetVideoByID(id int) (domain.Video, error)
	GetVideosByAuthorID(authorID int) ([]domain.Video, error)
	SearchVideosByTitle(title string) ([]domain.Video, error)
	IncrementVideoViews(id int) error
	UpdateVideo(video domain.Video) (domain.Video, error)
	CreateVideoQuality(quality domain.VideoQuality) (domain.VideoQuality, error)
	GetVideoQualities(videoID int) ([]domain.VideoQuality, error)
	DeleteVideo(id int) error
	ArchiveDeletedVideo(videoID int, archiveURL string) error

	CreateCategory(category domain.Category) (domain.Category, error)
	GetAllCategories() ([]domain.Category, error)
	GetCategoryByName(name string) (domain.Category, error)
	UpdateCategory(category domain.Category) (domain.Category, error)
	DeleteCategory(id int) error

	CreateDonation(donation domain.Donation) (domain.Donation, error)
	GetSentDonations(senderID int) ([]domain.Donation, error)
	GetReceivedDonations(receiverID int) ([]domain.Donation, error)
	GetUserDonations(userID int) ([]domain.Donation, error)

	LikeVideo(userID int, videoID int) error
	UnlikeVideo(userID int, videoID int) error
	GetVideoLikesCount(videoID int) (int, error)
	IsVideoLikedByUser(userID int, videoID int) (bool, error)

	SubscribeToUser(subscriberID int, authorID int) error
	UnsubscribeFromUser(subscriberID int, authorID int) error
	GetSubscribersCount(authorID int) (int, error)
	GetSubscriptionsCount(subscriberID int) (int, error)
	IsSubscribed(subscriberID int, authorID int) (bool, error)
	GetSubscribers(authorID int) ([]domain.User, error)
	GetSubscriptions(subscriberID int) ([]domain.User, error)

	CreateComment(comment domain.Comment) (domain.Comment, error)
	GetVideoComments(videoID int) ([]domain.Comment, error)
	GetCommentByID(commentID int) (domain.Comment, error)
	UpdateComment(comment domain.Comment) (domain.Comment, error)
	DeleteComment(commentID int) error

	GetChatByID(chatID int) (domain.Chat, error)
	GetUserChats(userID int) ([]domain.Chat, error)
	GetChatBetweenUsers(firstUserID int, secondUserID int) (domain.Chat, error)
	CreateChatRequest(request domain.ChatRequest) (domain.ChatRequest, error)
	GetChatRequestByID(requestID int) (domain.ChatRequest, error)
	GetIncomingChatRequests(userID int) ([]domain.ChatRequest, error)
	GetOutgoingChatRequests(userID int) ([]domain.ChatRequest, error)
	AcceptChatRequest(requestID int, firstUserID int, secondUserID int) (domain.ChatRequest, domain.Chat, error)
	RejectChatRequest(requestID int) (domain.ChatRequest, error)
	CreateChatMessage(message domain.ChatMessage) (domain.ChatMessage, error)
	GetChatMessages(chatID int) ([]domain.ChatMessage, error)

	CreatePlaylist(playlist domain.Playlist) (domain.Playlist, error)
	GetPlaylistByID(id int) (domain.Playlist, error)
	GetUserPlaylists(userID int) ([]domain.Playlist, error)
	UpdatePlaylist(playlist domain.Playlist) (domain.Playlist, error)
	DeletePlaylist(id int) error
	AddVideoToPlaylist(playlistID int, videoID int) error
	RemoveVideoFromPlaylist(playlistID int, videoID int) error
	GetPlaylistVideos(playlistID int) ([]domain.Video, error)
}
