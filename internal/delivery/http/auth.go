package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/delivery/http/dto"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/pkg"
)

func (h *Handler) SignUp(w http.ResponseWriter, r *http.Request) {
	var input dto.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	err := h.service.CreateUser(domain.User{
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CommonResponse{Message: "User created successfully"})
}

func (h *Handler) SignIn(w http.ResponseWriter, r *http.Request) {
	var input dto.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	userID, userRole, err := h.service.Authenticate(domain.User{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	accessToken, refreshToken, err := h.generateNewTokenPair(userID, userRole)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) RefreshTokenPair(w http.ResponseWriter, r *http.Request) {
	tokenString, err := h.extractTokenFromHeader(r, refreshTokenHeader)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, CommonError{Error: err.Error()})
		return
	}

	userID, isRefresh, userRole, err := pkg.ParseToken(tokenString)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, CommonError{Error: err.Error()})
		return
	}

	if !isRefresh {
		writeJSON(w, http.StatusUnauthorized, CommonError{Error: "use refresh token"})
		return
	}

	accessToken, refreshToken, err := h.generateNewTokenPair(userID, userRole)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	profile, err := h.service.GetUserProfile(userID, &userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, profile)
}
