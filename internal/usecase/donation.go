package usecase

import (
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func (uc *Usecase) CreateDonation(senderID int, donation domain.Donation) (domain.Donation, error) {
	donation.Message = strings.TrimSpace(donation.Message)

	if senderID <= 0 || donation.Amount <= 0 {
		return domain.Donation{}, errs.ErrInvalidFieldValue
	}

	if donation.VideoID != nil {
		if *donation.VideoID <= 0 {
			return domain.Donation{}, errs.ErrInvalidFieldValue
		}

		video, err := uc.GetVideoByID(*donation.VideoID)
		if err != nil {
			return domain.Donation{}, err
		}

		donation.ReceiverID = video.AuthorID
	}

	if donation.ReceiverID <= 0 {
		return domain.Donation{}, errs.ErrInvalidFieldValue
	}

	if senderID == donation.ReceiverID {
		return domain.Donation{}, errs.ErrCannotDonateToYourself
	}

	if _, err := uc.GetUserByID(donation.ReceiverID); err != nil {
		return domain.Donation{}, err
	}

	donation.SenderID = senderID
	return uc.repository.CreateDonation(donation)
}

func (uc *Usecase) GetSentDonations(userID int) ([]domain.Donation, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	return uc.repository.GetSentDonations(userID)
}

func (uc *Usecase) GetReceivedDonations(userID int) ([]domain.Donation, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	return uc.repository.GetReceivedDonations(userID)
}

func (uc *Usecase) GetUserDonations(userID int) ([]domain.Donation, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	if _, err := uc.GetUserByID(userID); err != nil {
		return nil, err
	}

	return uc.repository.GetUserDonations(userID)
}
