package contracts

import "github.com/Nonameipal/AnalogYouTube/internal/models/domain"

type RepositoryI interface {
	CreateUser(user domain.User) error
	GetUserByUsername(username string) (domain.User, error)
	GetUserByEmail(email string) (domain.User, error)
	GetUserByID(id int) (domain.User, error)
}
