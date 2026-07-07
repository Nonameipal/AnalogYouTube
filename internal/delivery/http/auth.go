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

// SignUp godoc
// @Summary Регистрация
// @Description Создает нового пользователя. Email можно не указывать.
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body dto.SignUpRequest true "Данные регистрации"
// @Success 201 {object} CommonResponse
// @Failure 400 {object} CommonError
// @Failure 422 {object} CommonError
// @Router /auth/sign-up [post]
// @Router /api/register [post]
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

// SignIn godoc
// @Summary Вход
// @Description Возвращает access и refresh токены.
// @Tags Auth
// @Accept json
// @Produce json
// @Param input body dto.SignInRequest true "Данные входа"
// @Success 200 {object} dto.TokenPairResponse
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Router /auth/sign-in [post]
// @Router /api/login [post]
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

// RefreshTokenPair godoc
// @Summary Обновить токены
// @Description Меняет refresh token на новую пару токенов.
// @Tags Auth
// @Produce json
// @Param Refresh-Token header string true "Refresh token"
// @Success 200 {object} dto.TokenPairResponse
// @Failure 401 {object} CommonError
// @Router /auth/refresh [get]
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

// Me godoc
// @Summary Мой профиль
// @Description Данные текущего пользователя.
// @Tags Profile
// @Produce json
// @Security BearerAuth
// @Success 200 {object} domain.UserProfile
// @Failure 401 {object} CommonError
// @Router /api/me [get]
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
