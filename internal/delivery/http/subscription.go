package httpdelivery

import (
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func (h *Handler) SubscribeToUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	authorID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = h.service.SubscribeToUser(userID, authorID); err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CommonResponse{Message: "Subscribed successfully"})
}

func (h *Handler) UnsubscribeFromUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	authorID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = h.service.UnsubscribeFromUser(userID, authorID); err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CommonResponse{Message: "Unsubscribed successfully"})
}

func (h *Handler) GetSubscribersCount(w http.ResponseWriter, r *http.Request) {
	authorID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	count, err := h.service.GetSubscribersCount(authorID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"subscribers_count": count})
}

func (h *Handler) GetSubscriptionsCount(w http.ResponseWriter, r *http.Request) {
	subscriberID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	count, err := h.service.GetSubscriptionsCount(subscriberID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"subscriptions_count": count})
}

func (h *Handler) IsSubscribedByMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	authorID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	subscribed, err := h.service.IsSubscribed(userID, authorID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"subscribed": subscribed})
}
