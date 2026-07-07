package usecase

import (
	"errors"
	"net/mail"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/utils"
)

func (uc *Usecase) CreateUser(user domain.User) error {
	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)
	user.Password = strings.TrimSpace(user.Password)

	if user.Username == "" || user.Password == "" {
		return errs.ErrInvalidFieldValue
	}

	normalizedEmail, err := normalizeEmail(user.Email)
	if err != nil {
		return err
	}
	user.Email = normalizedEmail

	_, err = uc.repository.GetUserByUsername(user.Username)
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return err
		}
	} else {
		return errs.ErrUsernameAlreadyExists
	}

	if user.Email != "" {
		_, err = uc.repository.GetUserByEmail(user.Email)
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

	return uc.repository.CreateUser(user)
}

func (uc *Usecase) Authenticate(user domain.User) (int, string, error) {
	user.Username = strings.TrimSpace(user.Username)
	user.Password = strings.TrimSpace(user.Password)

	if user.Username == "" || user.Password == "" {
		return 0, "", errs.ErrInvalidFieldValue
	}

	userFromDB, err := uc.repository.GetUserByUsername(user.Username)
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

func (uc *Usecase) GetUserByID(id int) (domain.User, error) {
	if id <= 0 {
		return domain.User{}, errs.ErrInvalidFieldValue
	}

	user, err := uc.repository.GetUserByID(id)
	if err != nil {
		if errors.Is(err, errs.ErrNotFound) {
			return domain.User{}, errs.ErrUserNotFound
		}
		return domain.User{}, err
	}

	return user, nil
}

func (uc *Usecase) UpdateUserProfile(userID int, user domain.User) (domain.User, error) {
	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)
	user.Password = strings.TrimSpace(user.Password)
	user.AvatarURL = strings.TrimSpace(user.AvatarURL)
	user.Description = strings.TrimSpace(user.Description)

	if userID <= 0 || user.Username == "" {
		return domain.User{}, errs.ErrInvalidFieldValue
	}

	normalizedEmail, err := normalizeEmail(user.Email)
	if err != nil {
		return domain.User{}, err
	}
	user.Email = normalizedEmail

	oldUser, err := uc.GetUserByID(userID)
	if err != nil {
		return domain.User{}, err
	}

	userByUsername, err := uc.repository.GetUserByUsername(user.Username)
	if err != nil {
		if !errors.Is(err, errs.ErrNotFound) {
			return domain.User{}, err
		}
	} else if userByUsername.ID != userID {
		return domain.User{}, errs.ErrUsernameAlreadyExists
	}

	if user.Email != "" {
		userByEmail, err := uc.repository.GetUserByEmail(user.Email)
		if err != nil {
			if !errors.Is(err, errs.ErrNotFound) {
				return domain.User{}, err
			}
		} else if userByEmail.ID != userID {
			return domain.User{}, errs.ErrEmailAlreadyExists
		}
	}

	user.ID = userID
	if user.AvatarURL == "" {
		user.AvatarURL = oldUser.AvatarURL
	}
	if user.Password == "" {
		user.Password = oldUser.Password
	} else {
		hashedPassword, err := utils.GenerateHash(user.Password)
		if err != nil {
			return domain.User{}, err
		}
		user.Password = hashedPassword
	}

	return uc.repository.UpdateUserProfile(user)
}

func (uc *Usecase) GetUserProfile(userID int, viewerID *int) (domain.UserProfile, error) {
	if userID <= 0 {
		return domain.UserProfile{}, errs.ErrInvalidFieldValue
	}

	user, err := uc.GetUserByID(userID)
	if err != nil {
		return domain.UserProfile{}, err
	}

	videos, err := uc.repository.GetVideosByAuthorID(userID)
	if err != nil {
		return domain.UserProfile{}, err
	}

	subscribersCount, err := uc.repository.GetSubscribersCount(userID)
	if err != nil {
		return domain.UserProfile{}, err
	}

	subscriptionsCount, err := uc.repository.GetSubscriptionsCount(userID)
	if err != nil {
		return domain.UserProfile{}, err
	}

	profile := domain.UserProfile{
		ID:                 user.ID,
		Username:           user.Username,
		Email:              user.Email,
		Role:               user.Role,
		AvatarURL:          user.AvatarURL,
		Description:        user.Description,
		SubscribersCount:   subscribersCount,
		SubscriptionsCount: subscriptionsCount,
		Videos:             videos,
		CreatedAt:          user.CreatedAt,
		UpdatedAt:          user.UpdatedAt,
	}

	if viewerID != nil && *viewerID == userID {
		subscribers, err := uc.repository.GetSubscribers(userID)
		if err != nil {
			return domain.UserProfile{}, err
		}

		subscriptions, err := uc.repository.GetSubscriptions(userID)
		if err != nil {
			return domain.UserProfile{}, err
		}

		profile.Subscribers = &subscribers
		profile.Subscriptions = &subscriptions
	}

	return profile, nil
}

func normalizeEmail(email string) (string, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return "", nil
	}
	if strings.ContainsAny(email, " \t\r\n") {
		return "", errs.ErrInvalidEmail
	}

	address, err := mail.ParseAddress(email)
	if err != nil || address.Address != email {
		return "", errs.ErrInvalidEmail
	}

	localPart, domainPart, ok := strings.Cut(email, "@")
	if !ok || localPart == "" || domainPart == "" || strings.Count(email, "@") != 1 {
		return "", errs.ErrInvalidEmail
	}
	if !isValidEmailDomain(domainPart) {
		return "", errs.ErrInvalidEmail
	}

	return email, nil
}

func isValidEmailDomain(domain string) bool {
	labels := strings.Split(domain, ".")
	if len(labels) < 2 {
		return false
	}

	for _, label := range labels {
		if label == "" || strings.HasPrefix(label, "-") || strings.HasSuffix(label, "-") {
			return false
		}
		for _, char := range label {
			isLetter := char >= 'a' && char <= 'z'
			isDigit := char >= '0' && char <= '9'
			if !isLetter && !isDigit && char != '-' {
				return false
			}
		}
	}

	topLevelDomain := labels[len(labels)-1]
	if len(topLevelDomain) < 2 {
		return false
	}
	for _, char := range topLevelDomain {
		if char < 'a' || char > 'z' {
			return false
		}
	}

	return true
}
