package contracts

import "github.com/Nonameipal/AnalogYouTube/internal/models/domain"

type ServiceI interface {
	CreateUser(user domain.User) error
	Authenticate(user domain.User) (int, string, error)
	GetUserByID(id int) (domain.User, error)
	UpdateUserProfile(userID int, user domain.User) (domain.User, error)

	CreateVideo(authorID int, video domain.Video) (domain.Video, error)
	GetAllVideos() ([]domain.Video, error)
	GetRecommendedVideos() ([]domain.Video, error)
	GetVideoByID(id int) (domain.Video, error)
	GetUserVideos(userID int) ([]domain.Video, error)
	IncrementVideoViews(id int) error
	UpdateVideo(userID int, userRole string, video domain.Video) (domain.Video, error)
	DeleteVideo(userID int, userRole string, videoID int) error

	CreateCategory(category domain.Category) (domain.Category, error)
	GetAllCategories() ([]domain.Category, error)
	GetCategoryByID(id int) (domain.Category, error)
	UpdateCategory(category domain.Category) (domain.Category, error)
	DeleteCategory(id int) error

	CreateDonation(senderID int, donation domain.Donation) (domain.Donation, error)
	GetSentDonations(userID int) ([]domain.Donation, error)
	GetReceivedDonations(userID int) ([]domain.Donation, error)
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

	CreateComment(userID int, videoID int, comment domain.Comment) (domain.Comment, error)
	GetVideoComments(videoID int) ([]domain.Comment, error)
	UpdateComment(userID int, userRole string, comment domain.Comment) (domain.Comment, error)
	DeleteComment(userID int, userRole string, commentID int) error
}
