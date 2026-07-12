package usecase

import (
	"errors"
	"testing"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/utils"
)

func TestNormalizeEmail(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr error
	}{
		{name: "empty email is allowed", input: " ", want: ""},
		{name: "lowercase and trim", input: " User@Gmail.COM ", want: "user@gmail.com"},
		{name: "missing at", input: "user.gmail.com", wantErr: errs.ErrInvalidEmail},
		{name: "spaces inside", input: "user name@gmail.com", wantErr: errs.ErrInvalidEmail},
		{name: "domain without dot", input: "user@gmail", wantErr: errs.ErrInvalidEmail},
		{name: "bad domain label", input: "user@-gmail.com", wantErr: errs.ErrInvalidEmail},
		{name: "bad top level domain", input: "user@gmail.c1", wantErr: errs.ErrInvalidEmail},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeEmail(tt.input)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if got != tt.want {
				t.Fatalf("expected email %q, got %q", tt.want, got)
			}
		})
	}
}

func TestCreateUser(t *testing.T) {
	t.Run("trims normalizes hashes and saves", func(t *testing.T) {
		repo := newFakeRepository()
		uc := newTestUsecase(repo)

		err := uc.CreateUser(domain.User{
			Username: "  neo  ",
			Email:    " Neo@Gmail.COM ",
			Password: "  secret  ",
		})
		if err != nil {
			t.Fatalf("CreateUser returned error: %v", err)
		}

		if repo.createdUser.Username != "neo" {
			t.Fatalf("expected trimmed username, got %q", repo.createdUser.Username)
		}
		if repo.createdUser.Email != "neo@gmail.com" {
			t.Fatalf("expected normalized email, got %q", repo.createdUser.Email)
		}
		if repo.createdUser.Role != domain.UserRole {
			t.Fatalf("expected role %q, got %q", domain.UserRole, repo.createdUser.Role)
		}
		if repo.createdUser.Password == "secret" {
			t.Fatal("password was saved without hashing")
		}
		if err := utils.CompareHash(repo.createdUser.Password, "secret"); err != nil {
			t.Fatalf("saved password hash does not match original password: %v", err)
		}
	})

	t.Run("rejects duplicate username", func(t *testing.T) {
		repo := newFakeRepository()
		repo.usersByUsername["neo"] = domain.User{ID: 1, Username: "neo"}
		uc := newTestUsecase(repo)

		err := uc.CreateUser(domain.User{Username: "neo", Password: "secret"})
		if !errors.Is(err, errs.ErrUsernameAlreadyExists) {
			t.Fatalf("expected ErrUsernameAlreadyExists, got %v", err)
		}
	})

	t.Run("rejects duplicate email", func(t *testing.T) {
		repo := newFakeRepository()
		repo.usersByEmail["neo@gmail.com"] = domain.User{ID: 1, Email: "neo@gmail.com"}
		uc := newTestUsecase(repo)

		err := uc.CreateUser(domain.User{Username: "neo", Email: "Neo@Gmail.COM", Password: "secret"})
		if !errors.Is(err, errs.ErrEmailAlreadyExists) {
			t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
		}
	})

	t.Run("rejects empty username or password", func(t *testing.T) {
		uc := newTestUsecase(newFakeRepository())

		err := uc.CreateUser(domain.User{Username: " ", Password: "secret"})
		if !errors.Is(err, errs.ErrInvalidFieldValue) {
			t.Fatalf("expected ErrInvalidFieldValue, got %v", err)
		}
	})
}

func TestAuthenticate(t *testing.T) {
	hash, err := utils.GenerateHash("secret")
	if err != nil {
		t.Fatalf("failed to prepare password hash: %v", err)
	}

	repo := newFakeRepository()
	repo.usersByUsername["neo"] = domain.User{ID: 7, Username: "neo", Password: hash, Role: domain.AdminRole}
	uc := newTestUsecase(repo)

	userID, role, err := uc.Authenticate(domain.User{Username: " neo ", Password: " secret "})
	if err != nil {
		t.Fatalf("Authenticate returned error: %v", err)
	}
	if userID != 7 || role != domain.AdminRole {
		t.Fatalf("expected user 7 with admin role, got user=%d role=%q", userID, role)
	}

	_, _, err = uc.Authenticate(domain.User{Username: "neo", Password: "bad"})
	if !errors.Is(err, errs.ErrIncorrectUsernameOrPassword) {
		t.Fatalf("expected ErrIncorrectUsernameOrPassword, got %v", err)
	}

	_, _, err = uc.Authenticate(domain.User{Username: "missing", Password: "secret"})
	if !errors.Is(err, errs.ErrUserNotFound) {
		t.Fatalf("expected ErrUserNotFound, got %v", err)
	}
}

func TestUpdateUserProfile(t *testing.T) {
	oldHash, err := utils.GenerateHash("old")
	if err != nil {
		t.Fatalf("failed to prepare password hash: %v", err)
	}

	t.Run("keeps old avatar and password when omitted", func(t *testing.T) {
		repo := newFakeRepository()
		repo.usersByID[1] = domain.User{ID: 1, Username: "old", Email: "old@gmail.com", Password: oldHash, AvatarURL: "/old.png"}
		uc := newTestUsecase(repo)

		user, err := uc.UpdateUserProfile(1, domain.User{
			Username:    " neo ",
			Email:       " Neo@Gmail.COM ",
			Description: " channel ",
		})
		if err != nil {
			t.Fatalf("UpdateUserProfile returned error: %v", err)
		}
		if user.AvatarURL != "/old.png" {
			t.Fatalf("expected old avatar, got %q", user.AvatarURL)
		}
		if user.Password != oldHash {
			t.Fatal("expected old password hash to be kept")
		}
		if user.Email != "neo@gmail.com" || user.Description != "channel" {
			t.Fatalf("profile was not normalized: %+v", user)
		}
	})

	t.Run("hashes new password", func(t *testing.T) {
		repo := newFakeRepository()
		repo.usersByID[1] = domain.User{ID: 1, Username: "old", Password: oldHash}
		uc := newTestUsecase(repo)

		user, err := uc.UpdateUserProfile(1, domain.User{Username: "neo", Password: "new"})
		if err != nil {
			t.Fatalf("UpdateUserProfile returned error: %v", err)
		}
		if user.Password == "new" || user.Password == oldHash {
			t.Fatal("new password was not hashed and saved")
		}
		if err := utils.CompareHash(user.Password, "new"); err != nil {
			t.Fatalf("new password hash does not match: %v", err)
		}
	})

	t.Run("rejects username owned by another user", func(t *testing.T) {
		repo := newFakeRepository()
		repo.usersByID[1] = domain.User{ID: 1, Username: "old"}
		repo.usersByUsername["taken"] = domain.User{ID: 2, Username: "taken"}
		uc := newTestUsecase(repo)

		_, err := uc.UpdateUserProfile(1, domain.User{Username: "taken"})
		if !errors.Is(err, errs.ErrUsernameAlreadyExists) {
			t.Fatalf("expected ErrUsernameAlreadyExists, got %v", err)
		}
	})
}

func TestGetUserProfile(t *testing.T) {
	repo := newFakeRepository()
	repo.usersByID[1] = domain.User{ID: 1, Username: "neo", Email: "neo@gmail.com", Role: domain.UserRole, Description: "channel"}
	repo.videosByID[10] = domain.Video{ID: 10, AuthorID: 1, Title: "first"}
	repo.subscribers = []domain.User{{ID: 2, Username: "trinity"}}
	repo.subscriptions = []domain.User{{ID: 3, Username: "morpheus"}}
	uc := newTestUsecase(repo)

	viewerID := 1
	profile, err := uc.GetUserProfile(1, &viewerID)
	if err != nil {
		t.Fatalf("GetUserProfile returned error: %v", err)
	}
	if profile.Username != "neo" || profile.SubscribersCount != 4 || profile.SubscriptionsCount != 2 {
		t.Fatalf("unexpected public profile: %+v", profile)
	}
	if profile.Subscribers == nil || profile.Subscriptions == nil {
		t.Fatal("owner should see subscribers and subscriptions lists")
	}

	otherViewerID := 2
	profile, err = uc.GetUserProfile(1, &otherViewerID)
	if err != nil {
		t.Fatalf("GetUserProfile returned error: %v", err)
	}
	if profile.Subscribers != nil || profile.Subscriptions != nil {
		t.Fatal("other users should only see counters")
	}
}
