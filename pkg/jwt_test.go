package pkg

import (
	"strings"
	"testing"

	"github.com/Nonameipal/AnalogYouTube/internal/configs"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
)

func TestGenerateAndParseToken(t *testing.T) {
	configs.AppSettings.AuthParams.JwtSecret = "test-secret"

	token, err := GenerateToken(7, 15, domain.AdminRole, false)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	userID, isRefresh, role, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
	}
	if userID != 7 || isRefresh || role != domain.AdminRole {
		t.Fatalf("unexpected token data: userID=%d isRefresh=%v role=%q", userID, isRefresh, role)
	}
}

func TestGenerateAndParseRefreshToken(t *testing.T) {
	configs.AppSettings.AuthParams.JwtSecret = "test-secret"

	token, err := GenerateToken(9, 30, domain.UserRole, true)
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}

	userID, isRefresh, role, err := ParseToken(token)
	if err != nil {
		t.Fatalf("ParseToken returned error: %v", err)
	}
	if userID != 9 || !isRefresh || role != domain.UserRole {
		t.Fatalf("unexpected token data: userID=%d isRefresh=%v role=%q", userID, isRefresh, role)
	}
}

func TestParseTokenRejectsInvalidToken(t *testing.T) {
	configs.AppSettings.AuthParams.JwtSecret = "test-secret"

	_, _, _, err := ParseToken("not-a-token")
	if err == nil || !strings.Contains(err.Error(), "token") {
		t.Fatalf("expected token parsing error, got %v", err)
	}
}
