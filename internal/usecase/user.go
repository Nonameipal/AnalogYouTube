package usecase

import (
	"errors"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/utils"
)

func (s *Service) CreateUser(user domain.User) error {
	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)
	user.Password = strings.TrimSpace(user.Password)

	if user.Username == "" || user.Password == "" {
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
	if id <= 0 {
		return domain.User{}, errs.ErrInvalidFieldValue
	}

	user, err := s.repository.GetUserByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.User{}, errs.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (s *Service) UpdateUserProfile(userID int, user domain.User) (domain.User, error) {
	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)
	user.AvatarURL = strings.TrimSpace(user.AvatarURL)
	user.Description = strings.TrimSpace(user.Description)

	if userID <= 0 || user.Username == "" {
		return domain.User{}, errs.ErrInvalidFieldValue
	}

	if _, err := s.GetUserByID(userID); err != nil {
		return domain.User{}, err
	}

	userByUsername, err := s.repository.GetUserByUsername(user.Username)
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return domain.User{}, err
		}
	} else if userByUsername.ID != userID {
		return domain.User{}, errs.ErrUsernameAlreadyExists
	}

	if user.Email != "" {
		userByEmail, err := s.repository.GetUserByEmail(user.Email)
		if err != nil {
			if !errors.Is(err, errs.ErrNotFound) {
				return domain.User{}, err
			}
		} else if userByEmail.ID != userID {
			return domain.User{}, errs.ErrEmailAlreadyExists
		}
	}

	user.ID = userID
	return s.repository.UpdateUserProfile(user)
}

func (s *Service) GetUserProfile(userID int, viewerID *int) (domain.UserProfile, error) {
	if userID <= 0 {
		return domain.UserProfile{}, errs.ErrInvalidFieldValue
	}

	user, err := s.GetUserByID(userID)
	if err != nil {
		return domain.UserProfile{}, err
	}

	videos, err := s.repository.GetVideosByAuthorID(userID)
	if err != nil {
		return domain.UserProfile{}, err
	}

	subscribersCount, err := s.repository.GetSubscribersCount(userID)
	if err != nil {
		return domain.UserProfile{}, err
	}

	subscriptionsCount, err := s.repository.GetSubscriptionsCount(userID)
	if err != nil {
		return domain.UserProfile{}, err
	}

	profile := domain.UserProfile{
		ID: user.ID,
		Username: user.Username,
		Email: user.Email,
		Role: user.Role,
		AvatarURL: user.AvatarURL,
		Description: user.Description,
		SubscribersCount: subscribersCount,
		SubscriptionsCount: subscriptionsCount,
		Videos: videos,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	if viewerID != nil && *viewerID == userID {
		subscribers, err := s.repository.GetSubscribers(userID)
		if err != nil {
			return domain.UserProfile{}, err
		}

		subscriptions, err := s.repository.GetSubscriptions(userID)
		if err != nil {
			return domain.UserProfile{}, err
		}

		profile.Subscribers = &subscribers
		profile.Subscriptions = &subscriptions
	}

	return profile, nil
}

func (s *Service) UpdateUserAvatar(userID int, avatarURL string) (domain.User, error) {
	avatarURL = strings.TrimSpace(avatarURL)

	if userID <= 0 || avatarURL == "" {
		return domain.User{}, errs.ErrInvalidFieldValue
	}

	if _, err := s.GetUserByID(userID); err != nil {
		return domain.User{}, err
	}

	user, err := s.repository.UpdateUserAvatar(userID, avatarURL)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.User{}, errs.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}
