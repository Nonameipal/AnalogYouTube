package usecase

import "github.com/Nonameipal/AnalogYouTube/internal/errs"

func (uc *Usecase) SubscribeToUser(subscriberID int, authorID int) error {
	if subscriberID <= 0 || authorID <= 0 {
		return errs.ErrInvalidFieldValue
	}
	if subscriberID == authorID {
		return errs.ErrCannotSubscribeToYourself
	}

	if _, err := uc.GetUserByID(authorID); err != nil {
		return err
	}

	return uc.repository.SubscribeToUser(subscriberID, authorID)
}

func (uc *Usecase) UnsubscribeFromUser(subscriberID int, authorID int) error {
	if subscriberID <= 0 || authorID <= 0 {
		return errs.ErrInvalidFieldValue
	}
	if subscriberID == authorID {
		return errs.ErrCannotSubscribeToYourself
	}

	if _, err := uc.GetUserByID(authorID); err != nil {
		return err
	}

	return uc.repository.UnsubscribeFromUser(subscriberID, authorID)
}

func (uc *Usecase) GetSubscribersCount(authorID int) (int, error) {
	if authorID <= 0 {
		return 0, errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetUserByID(authorID); err != nil {
		return 0, err
	}

	return uc.repository.GetSubscribersCount(authorID)
}

func (uc *Usecase) GetSubscriptionsCount(subscriberID int) (int, error) {
	if subscriberID <= 0 {
		return 0, errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetUserByID(subscriberID); err != nil {
		return 0, err
	}

	return uc.repository.GetSubscriptionsCount(subscriberID)
}

func (uc *Usecase) IsSubscribed(subscriberID int, authorID int) (bool, error) {
	if subscriberID <= 0 || authorID <= 0 {
		return false, errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetUserByID(authorID); err != nil {
		return false, err
	}

	return uc.repository.IsSubscribed(subscriberID, authorID)
}
