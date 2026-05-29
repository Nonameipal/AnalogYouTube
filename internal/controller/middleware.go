package controller

import (
	"context"
	"net/http"

	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/internal/models/domain"
	"github.com/Nonameipal/AnalogYouTube/pkg"
)

type contextKey string

const (
	userIDContextKey   contextKey = "user_id"
	userRoleContextKey contextKey = "user_role"
)

func (ctrl *Controller) checkUserAuthentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tokenString, err := ctrl.extractTokenFromHeader(r, authorizationHeader)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, CommonError{Error: err.Error()})
			return
		}

		userID, isRefresh, userRole, err := pkg.ParseToken(tokenString)
		if err != nil {
			writeJSON(w, http.StatusUnauthorized, CommonError{Error: err.Error()})
			return
		}

		if isRefresh {
			writeJSON(w, http.StatusUnauthorized, CommonError{Error: "use access token, not refresh token"})
			return
		}

		ctx := context.WithValue(r.Context(), userIDContextKey, userID)
		ctx = context.WithValue(ctx, userRoleContextKey, userRole)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (ctrl *Controller) checkIsAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role, ok := r.Context().Value(userRoleContextKey).(string)
		if !ok || role != domain.AdminRole {
			ctrl.handleError(w, errs.ErrAccessDenied)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func getUserIDFromContext(r *http.Request) (int, bool) {
	userID, ok := r.Context().Value(userIDContextKey).(int)
	return userID, ok
}
