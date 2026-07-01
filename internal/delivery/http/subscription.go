package httpdelivery

import (
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

func (ctrl *Controller) SubscribeToUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	authorID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	if err = ctrl.service.SubscribeToUser(userID, authorID); err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CommonResponse{Message: "Subscribed successfully"})
}

func (ctrl *Controller) UnsubscribeFromUser(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	authorID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	if err = ctrl.service.UnsubscribeFromUser(userID, authorID); err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, CommonResponse{Message: "Unsubscribed successfully"})
}

func (ctrl *Controller) GetSubscribersCount(w http.ResponseWriter, r *http.Request) {
	authorID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	count, err := ctrl.service.GetSubscribersCount(authorID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"subscribers_count": count})
}

func (ctrl *Controller) GetSubscriptionsCount(w http.ResponseWriter, r *http.Request) {
	subscriberID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	count, err := ctrl.service.GetSubscriptionsCount(subscriberID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]int{"subscriptions_count": count})
}

func (ctrl *Controller) IsSubscribedByMe(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	authorID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	subscribed, err := ctrl.service.IsSubscribed(userID, authorID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, map[string]bool{"subscribed": subscribed})
}
