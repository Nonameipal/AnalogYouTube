package usecase

import (
	"errors"
	"testing"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func TestDonationUsecase(t *testing.T) {
	t.Run("direct donation trims message and saves sender", func(t *testing.T) {
		repo := newFakeRepository()
		repo.usersByID[2] = domain.User{ID: 2, Username: "receiver"}
		uc := newTestUsecase(repo)

		donation, err := uc.CreateDonation(1, domain.Donation{ReceiverID: 2, Amount: 10.5, Message: " thanks "})
		if err != nil {
			t.Fatalf("CreateDonation returned error: %v", err)
		}
		if donation.SenderID != 1 || donation.ReceiverID != 2 || donation.Message != "thanks" {
			t.Fatalf("unexpected donation: %+v", donation)
		}
	})

	t.Run("video donation uses video author as receiver", func(t *testing.T) {
		repo := newFakeRepository()
		repo.usersByID[7] = domain.User{ID: 7, Username: "author"}
		repo.videosByID[3] = domain.Video{ID: 3, AuthorID: 7}
		uc := newTestUsecase(repo)
		videoID := 3

		donation, err := uc.CreateDonation(1, domain.Donation{VideoID: &videoID, Amount: 5})
		if err != nil {
			t.Fatalf("CreateDonation returned error: %v", err)
		}
		if donation.ReceiverID != 7 {
			t.Fatalf("expected video author as receiver, got %d", donation.ReceiverID)
		}
	})

	t.Run("rejects self donation and invalid amount", func(t *testing.T) {
		repo := newFakeRepository()
		repo.usersByID[1] = domain.User{ID: 1, Username: "sender"}
		uc := newTestUsecase(repo)

		_, err := uc.CreateDonation(1, domain.Donation{ReceiverID: 1, Amount: 10})
		if !errors.Is(err, errs.ErrCannotDonateToYourself) {
			t.Fatalf("expected ErrCannotDonateToYourself, got %v", err)
		}

		_, err = uc.CreateDonation(1, domain.Donation{ReceiverID: 2, Amount: 0})
		if !errors.Is(err, errs.ErrInvalidFieldValue) {
			t.Fatalf("expected ErrInvalidFieldValue, got %v", err)
		}
	})
}
