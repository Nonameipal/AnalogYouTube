package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/models/dto"
	"github.com/Nonameipal/AnalogYouTube/pkg"
)

func (ctrl *Controller) SignUp(w http.ResponseWriter, r *http.Request) {
	var input dto.SignUpRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	err := ctrl.service.CreateUser(domain.User{
		FullName: input.FullName,
		Username: input.Username,
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, CommonResponse{Message: "User created successfully"})
}

func (ctrl *Controller) SignIn(w http.ResponseWriter, r *http.Request) {
	var input dto.SignInRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	userID, userRole, err := ctrl.service.Authenticate(domain.User{
		Username: input.Username,
		Password: input.Password,
	})
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	accessToken, refreshToken, err := ctrl.generateNewTokenPair(userID, userRole)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (ctrl *Controller) RefreshTokenPair(w http.ResponseWriter, r *http.Request) {
	tokenString, err := ctrl.extractTokenFromHeader(r, refreshTokenHeader)
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

	accessToken, refreshToken, err := ctrl.generateNewTokenPair(userID, userRole)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, dto.TokenPairResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}

func (ctrl *Controller) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	user, err := ctrl.service.GetUserByID(userID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}
