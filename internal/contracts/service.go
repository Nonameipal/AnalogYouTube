package contracts

import "github.com/Nonameipal/AnalogYouTube/internal/models/domain"

type ServiceI interface {
	CreateUser(user domain.User) error
	Authenticate(user domain.User) (int, string, error)
	GetUserByID(id int) (domain.User, error)

	CreateVideo(authorID int, video domain.Video) (domain.Video, error)
	GetAllVideos() ([]domain.Video, error)
	GetVideoByID(id int) (domain.Video, error)
	UpdateVideo(userID int, userRole string, video domain.Video) (domain.Video, error)
	DeleteVideo(userID int, userRole string, videoID int) error
}
