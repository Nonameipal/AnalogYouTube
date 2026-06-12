package controller

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/models/dto"
)

func (ctrl *Controller) CreateDonation(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	var input dto.CreateDonationRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	donation, err := ctrl.service.CreateDonation(userID, domain.Donation{
		ReceiverID: input.ReceiverID,
		VideoID:    input.VideoID,
		Amount:     input.Amount,
		Message:    input.Message,
	})
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, donation)
}

func (ctrl *Controller) GetSentDonations(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	donations, err := ctrl.service.GetSentDonations(userID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, donations)
}

func (ctrl *Controller) GetReceivedDonations(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	donations, err := ctrl.service.GetReceivedDonations(userID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, donations)
}

func (ctrl *Controller) GetUserDonations(w http.ResponseWriter, r *http.Request) {
	userID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	donations, err := ctrl.service.GetUserDonations(userID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, donations)
}
