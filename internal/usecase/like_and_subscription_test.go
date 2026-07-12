package usecase

import (
	"errors"
	"testing"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func TestLikeAndSubscriptionUsecase(t *testing.T) {
	t.Run("like checks video then saves", func(t *testing.T) {
		repo := newFakeRepository()
		repo.videosByID[3] = domain.Video{ID: 3}
		uc := newTestUsecase(repo)

		if err := uc.LikeVideo(1, 3); err != nil {
			t.Fatalf("LikeVideo returned error: %v", err)
		}
		if repo.likedUserID != 1 || repo.likedVideoID != 3 {
			t.Fatalf("like was not saved correctly")
		}
	})

	t.Run("subscribe rejects self subscription", func(t *testing.T) {
		uc := newTestUsecase(newFakeRepository())

		err := uc.SubscribeToUser(1, 1)
		if !errors.Is(err, errs.ErrCannotSubscribeToYourself) {
			t.Fatalf("expected ErrCannotSubscribeToYourself, got %v", err)
		}
	})

	t.Run("subscribe checks author and saves", func(t *testing.T) {
		repo := newFakeRepository()
		repo.usersByID[2] = domain.User{ID: 2, Username: "author"}
		uc := newTestUsecase(repo)

		err := uc.SubscribeToUser(1, 2)
		if err != nil {
			t.Fatalf("SubscribeToUser returned error: %v", err)
		}
		if repo.subscribedUserID != 1 || repo.subscribedAuthorID != 2 {
			t.Fatalf("subscription was not saved")
		}
	})
}
