package service

import "github.com/Nonameipal/AnalogYouTube/internal/errs"

func (s *Service) SubscribeToUser(subscriberID int, authorID int) error {
	if subscriberID <= 0 || authorID <= 0 {
		return errs.ErrInvalidFieldValue
	}
	if subscriberID == authorID {
		return errs.ErrCannotSubscribeToYourself
	}

	if _, err := s.GetUserByID(authorID); err != nil {
		return err
	}

	return s.repository.SubscribeToUser(subscriberID, authorID)
}

func (s *Service) UnsubscribeFromUser(subscriberID int, authorID int) error {
	if subscriberID <= 0 || authorID <= 0 {
		return errs.ErrInvalidFieldValue
	}
	if subscriberID == authorID {
		return errs.ErrCannotSubscribeToYourself
	}

	if _, err := s.GetUserByID(authorID); err != nil {
		return err
	}

	return s.repository.UnsubscribeFromUser(subscriberID, authorID)
}

func (s *Service) GetSubscribersCount(authorID int) (int, error) {
	if authorID <= 0 {
		return 0, errs.ErrInvalidFieldValue
	}

	if _, err := s.GetUserByID(authorID); err != nil {
		return 0, err
	}

	return s.repository.GetSubscribersCount(authorID)
}

func (s *Service) GetSubscriptionsCount(subscriberID int) (int, error) {
	if subscriberID <= 0 {
		return 0, errs.ErrInvalidFieldValue
	}

	if _, err := s.GetUserByID(subscriberID); err != nil {
		return 0, err
	}

	return s.repository.GetSubscriptionsCount(subscriberID)
}

func (s *Service) IsSubscribed(subscriberID int, authorID int) (bool, error) {
	if subscriberID <= 0 || authorID <= 0 {
		return false, errs.ErrInvalidFieldValue
	}

	if _, err := s.GetUserByID(authorID); err != nil {
		return false, err
	}

	return s.repository.IsSubscribed(subscriberID, authorID)
}
