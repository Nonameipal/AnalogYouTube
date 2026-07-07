package httpdelivery

import (
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
)

// SubscribeToUser godoc
// @Summary Подписаться
// @Description Подписка на автора.
// @Tags Subscriptions
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID пользователя, на которого подписываемся"
// @Success 200 {object} CommonResponse
// @Failure 400 {object} CommonError
// @Failure 401 {object} CommonError
// @Router /api/users/{id}/subscribe [post]
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

// UnsubscribeFromUser godoc
// @Summary Отписаться
// @Description Отмена подписки.
// @Tags Subscriptions
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Success 200 {object} CommonResponse
// @Failure 401 {object} CommonError
// @Router /api/users/{id}/subscribe [delete]
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

// GetSubscribersCount godoc
// @Summary Подписчики
// @Description Сколько людей подписано на автора.
// @Tags Subscriptions
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]int
// @Failure 400 {object} CommonError
// @Router /api/users/{id}/subscribers/count [get]
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

// GetSubscriptionsCount godoc
// @Summary Подписки
// @Description Сколько подписок у пользователя.
// @Tags Subscriptions
// @Produce json
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]int
// @Failure 400 {object} CommonError
// @Router /api/users/{id}/subscriptions/count [get]
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

// IsSubscribedByMe godoc
// @Summary Моя подписка
// @Description Показывает, подписан ли я на автора.
// @Tags Subscriptions
// @Produce json
// @Security BearerAuth
// @Param id path int true "ID пользователя"
// @Success 200 {object} map[string]bool
// @Failure 401 {object} CommonError
// @Router /api/users/{id}/subscribed [get]
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
