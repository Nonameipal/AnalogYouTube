package contracts

import "github.com/Nonameipal/AnalogYouTube/internal/models/domain"

type RepositoryI interface {
	CreateUser(user domain.User) error
	GetUserByUsername(username string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	GetUserByID(id int) (domain.User, error)

	CreateVideo(video domain.Video) (domain.Video, error)
	GetAllVideos() ([]domain.Video, error)
	GetVideoByID(id int) (domain.Video, error)
	UpdateVideo(video domain.Video) (domain.Video, error)
	DeleteVideo(id int) error
}
