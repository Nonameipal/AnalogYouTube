package contracts

import "github.com/Nonameipal/AnalogYouTube/internal/models/domain"

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
	IncrementVideoViews(id int) error
	UpdateVideo(video domain.Video) (domain.Video, error)
	DeleteVideo(id int) error

	CreateCategory(category domain.Category) (domain.Category, error)
	GetAllCategories() ([]domain.Category, error)
	GetCategoryByID(id int) (domain.Category, error)
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

	CreateComment(comment domain.Comment) (domain.Comment, error)
	GetVideoComments(videoID int) ([]domain.Comment, error)
	GetCommentByID(commentID int) (domain.Comment, error)
	UpdateComment(comment domain.Comment) (domain.Comment, error)
	DeleteComment(commentID int) error
}
