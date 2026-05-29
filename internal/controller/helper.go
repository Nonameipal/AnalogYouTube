package controller

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Nonameipal/AnalogYouTube/internal/configs"
	"github.com/Nonameipal/AnalogYouTube/pkg"
)

const (
	authorizationHeader = "Authorization"
	refreshTokenHeader  = "X-Refresh-Token"
)

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (ctrl *Controller) extractTokenFromHeader(r *http.Request, headerKey string) (string, error) {
	header := r.Header.Get(headerKey)
	if header == "" {
		return "", errors.New("empty authorization header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		return "", errors.New("invalid authorization header")
	}

	if !strings.EqualFold(headerParts[0], "Bearer") {
		return "", errors.New("authorization header must start with Bearer")
	}

	if headerParts[1] == "" {
		return "", errors.New("empty token")
	}

	return headerParts[1], nil
}

func (ctrl *Controller) generateNewTokenPair(userID int, userRole string) (string, string, error) {
	accessToken, err := pkg.GenerateToken(
		userID,
		configs.AppSettings.AuthParams.AccessTokenTtlMinutes,
		userRole,
		false,
	)
	if err != nil {
		return "", "", err
	}

	refreshToken, err := pkg.GenerateToken(
		userID,
		configs.AppSettings.AuthParams.RefreshTokenTtlDays,
		userRole,
		true,
	)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}
