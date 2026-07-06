package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/delivery/http/dto"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

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
