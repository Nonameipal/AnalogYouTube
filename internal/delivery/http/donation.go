package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/delivery/http/dto"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

// CreateDonation godoc
// @Summary Создать донат
// @Description Сохраняет донат без настоящей оплаты.
// @Tags Donations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param input body dto.CreateDonationRequest true "Данные доната"
// @Success 201 {object} domain.Donation
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Router /api/donations [post]
func (h *Handler) CreateDonation(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	var input dto.CreateDonationRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	donation, err := h.service.CreateDonation(userID, domain.Donation{
		ReceiverID: input.TargetUserID(),
		VideoID:    input.VideoID,
		Amount:     input.Amount,
		Message:    input.Message,
	})
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, donation)
}

// GetSentDonations godoc
// @Summary Мои отправленные донаты
// @Description Донаты, которые отправили вы.
// @Tags Donations
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Donation
// @Failure 401 {object} CommonError
// @Router /api/donations/sent [get]
func (h *Handler) GetSentDonations(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	donations, err := h.service.GetSentDonations(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, donations)
}

// GetReceivedDonations godoc
// @Summary Мои полученные донаты
// @Description Донаты, которые получили вы.
// @Tags Donations
// @Produce json
// @Security BearerAuth
// @Success 200 {array} domain.Donation
// @Failure 401 {object} CommonError
// @Router /api/donations/received [get]
func (h *Handler) GetReceivedDonations(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	donations, err := h.service.GetReceivedDonations(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, donations)
}

// GetUserDonations godoc
// @Summary Донаты пользователя
// @Description Публичный список полученных донатов.
// @Tags Donations
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {array} domain.Donation
// @Failure 400 {object} CommonError
// @Router /api/users/{id}/donations [get]
func (h *Handler) GetUserDonations(w http.ResponseWriter, r *http.Request) {
	userID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	donations, err := h.service.GetUserDonations(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, donations)
}
