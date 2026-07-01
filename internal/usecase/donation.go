package usecase

import (
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func (s *Service) CreateDonation(senderID int, donation domain.Donation) (domain.Donation, error) {
	donation.Message = strings.TrimSpace(donation.Message)

	if senderID <= 0 || donation.ReceiverID <= 0 || donation.Amount <= 0 {
		return domain.Donation{}, errs.ErrInvalidFieldValue
	}

	if senderID == donation.ReceiverID {
		return domain.Donation{}, errs.ErrCannotDonateToYourself
	}

	if _, err := s.GetUserByID(donation.ReceiverID); err != nil {
		return domain.Donation{}, err
	}

	if donation.VideoID != nil {
		if *donation.VideoID <= 0 {
			return domain.Donation{}, errs.ErrInvalidFieldValue
		}

		if _, err := s.GetVideoByID(*donation.VideoID); err != nil {
			return domain.Donation{}, err
		}
	}

	donation.SenderID = senderID
	return s.repository.CreateDonation(donation)
}

func (s *Service) GetSentDonations(userID int) ([]domain.Donation, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	return s.repository.GetSentDonations(userID)
}

func (s *Service) GetReceivedDonations(userID int) ([]domain.Donation, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	return s.repository.GetReceivedDonations(userID)
}

func (s *Service) GetUserDonations(userID int) ([]domain.Donation, error) {
	if userID <= 0 {
		return nil, errs.ErrInvalidFieldValue
	}

	if _, err := s.GetUserByID(userID); err != nil {
		return nil, err
	}

	return s.repository.GetUserDonations(userID)
}
