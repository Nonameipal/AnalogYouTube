package service

import (
	"errors"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
	"github.com/Nonameipal/AnalogYouTube/utils"
)

func (s *Service) CreateUser(user domain.User) error {
	user.FullName = strings.TrimSpace(user.FullName)
	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)
	user.Password = strings.TrimSpace(user.Password)

	if user.FullName == "" || user.Username == "" || user.Password == "" {
		return errs.ErrInvalidFieldValue
	}

	_, err := s.repository.GetUserByUsername(user.Username)
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return err
		}
	} else {
		return errs.ErrUsernameAlreadyExists
	}

	if user.Email != "" {
		_, err = s.repository.GetUserByEmail(user.Email)
		if err != nil {
			if !errors.Is(err, errs.ErrNotFound) {
				return err
			}
		} else {
			return errs.ErrEmailAlreadyExists
		}
	}

	hashedPassword, err := utils.GenerateHash(user.Password)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.Role = domain.UserRole

	return s.repository.CreateUser(user)
}

func (s *Service) Authenticate(user domain.User) (int, string, error) {
	user.Username = strings.TrimSpace(user.Username)
	user.Password = strings.TrimSpace(user.Password)

	if user.Username == "" || user.Password == "" {
		return 0, "", errs.ErrInvalidFieldValue
	}

	userFromDB, err := s.repository.GetUserByUsername(user.Username)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return 0, "", errs.ErrUserNotFound
		}
		return 0, "", err
	}

	if err := utils.CompareHash(userFromDB.Password, user.Password); err != nil {
		return 0, "", errs.ErrIncorrectUsernameOrPassword
	}

	return userFromDB.ID, userFromDB.Role, nil
}

func (s *Service) GetUserByID(id int) (domain.User, error) {
	user, err := s.repository.GetUserByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.User{}, errs.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}
