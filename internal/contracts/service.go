package contracts

import "github.com/Nonameipal/AnalogYouTube/internal/models/domain"

type ServiceI interface {
	CreateUser(user domain.User) error
	Authenticate(user domain.User) (int, string, error)
	GetUserByID(id int) (domain.User, error)
}
