package httpdelivery

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/Nonameipal/AnalogYouTube/internal/configs"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/pkg"
	"github.com/gorilla/mux"
)

func TestWriteJSON(t *testing.T) {
	w := httptest.NewRecorder()

	writeJSON(w, http.StatusCreated, CommonResponse{Message: "created"})

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
	if contentType := w.Header().Get("Content-Type"); contentType != "application/json" {
		t.Fatalf("expected json content type, got %q", contentType)
	}

	var response CommonResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if response.Message != "created" {
		t.Fatalf("unexpected response: %+v", response)
	}
}

func TestHandleErrorStatusCodes(t *testing.T) {
	h := &Handler{}
	tests := []struct {
		name   string
		err    error
		status int
	}{
		{name: "not found", err: errs.ErrUserNotFound, status: http.StatusNotFound},
		{name: "bad request", err: errs.ErrInvalidFieldValue, status: http.StatusBadRequest},
		{name: "unauthorized", err: errs.ErrInvalidToken, status: http.StatusUnauthorized},
		{name: "forbidden", err: errs.ErrAccessDenied, status: http.StatusForbidden},
		{name: "unprocessable", err: errs.ErrEmailAlreadyExists, status: http.StatusUnprocessableEntity},
		{name: "internal", err: errors.New("boom"), status: http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			h.handleError(w, tt.err)
			if w.Code != tt.status {
				t.Fatalf("expected status %d, got %d", tt.status, w.Code)
			}
		})
	}
}

func TestExtractTokenFromHeader(t *testing.T) {
	h := &Handler{}
	tests := []struct {
		name    string
		header  string
		want    string
		wantErr bool
	}{
		{name: "valid bearer", header: "Bearer token", want: "token"},
		{name: "case insensitive bearer", header: "bearer token", want: "token"},
		{name: "empty", header: "", wantErr: true},
		{name: "bad format", header: "Bearer", wantErr: true},
		{name: "bad scheme", header: "Token token", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.header != "" {
				r.Header.Set(authorizationHeader, tt.header)
			}

			got, err := h.extractTokenFromHeader(r, authorizationHeader)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != tt.want {
				t.Fatalf("expected token %q, got %q", tt.want, got)
			}
		})
	}
}

func TestGenerateNewTokenPair(t *testing.T) {
	configs.AppSettings.AuthParams.JwtSecret = "test-secret"
	configs.AppSettings.AuthParams.AccessTokenTtlMinutes = 15
	configs.AppSettings.AuthParams.RefreshTokenTtlDays = 30
	h := &Handler{}

	accessToken, refreshToken, err := h.generateNewTokenPair(5, domain.UserRole)
	if err != nil {
		t.Fatalf("generateNewTokenPair returned error: %v", err)
	}

	userID, isRefresh, role, err := pkg.ParseToken(accessToken)
	if err != nil {
		t.Fatalf("failed to parse access token: %v", err)
	}
	if userID != 5 || isRefresh || role != domain.UserRole {
		t.Fatalf("unexpected access token data")
	}

	userID, isRefresh, role, err = pkg.ParseToken(refreshToken)
	if err != nil {
		t.Fatalf("failed to parse refresh token: %v", err)
	}
	if userID != 5 || !isRefresh || role != domain.UserRole {
		t.Fatalf("unexpected refresh token data")
	}
}

func TestCheckUserAuthentication(t *testing.T) {
	configs.AppSettings.AuthParams.JwtSecret = "test-secret"
	h := &Handler{}

	accessToken, err := pkg.GenerateToken(11, 15, domain.AdminRole, false)
	if err != nil {
		t.Fatalf("failed to generate access token: %v", err)
	}

	nextCalled := false
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		userID, ok := getUserIDFromContext(r)
		if !ok || userID != 11 {
			t.Fatalf("unexpected user id in context: %d ok=%v", userID, ok)
		}
		role, ok := getUserRoleFromContext(r)
		if !ok || role != domain.AdminRole {
			t.Fatalf("unexpected role in context: %q ok=%v", role, ok)
		}
		writeJSON(w, http.StatusOK, CommonResponse{Message: "ok"})
	})

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(authorizationHeader, "Bearer "+accessToken)
	w := httptest.NewRecorder()

	h.checkUserAuthentication(next).ServeHTTP(w, r)
	if !nextCalled {
		t.Fatal("next handler was not called")
	}
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	refreshToken, err := pkg.GenerateToken(11, 30, domain.AdminRole, true)
	if err != nil {
		t.Fatalf("failed to generate refresh token: %v", err)
	}
	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set(authorizationHeader, "Bearer "+refreshToken)
	w = httptest.NewRecorder()

	h.checkUserAuthentication(next).ServeHTTP(w, r)
	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected refresh token to be rejected, got %d", w.Code)
	}
}

func TestCheckIsAdmin(t *testing.T) {
	h := &Handler{}
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, http.StatusOK, CommonResponse{Message: "admin"})
	})

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), userRoleContextKey, domain.AdminRole))
	w := httptest.NewRecorder()
	h.checkIsAdmin(next).ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected admin request to pass, got %d", w.Code)
	}

	r = httptest.NewRequest(http.MethodGet, "/", nil)
	r = r.WithContext(context.WithValue(r.Context(), userRoleContextKey, domain.UserRole))
	w = httptest.NewRecorder()
	h.checkIsAdmin(next).ServeHTTP(w, r)
	if w.Code != http.StatusForbidden {
		t.Fatalf("expected user request to be forbidden, got %d", w.Code)
	}
}

func TestGetIDFromRequest(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/items/12", nil)
	r = mux.SetURLVars(r, map[string]string{"id": "12"})

	id, err := getIDFromRequest(r, "id")
	if err != nil {
		t.Fatalf("getIDFromRequest returned error: %v", err)
	}
	if id != 12 {
		t.Fatalf("expected id 12, got %d", id)
	}

	r = mux.SetURLVars(r, map[string]string{"id": "bad"})
	if _, err := getIDFromRequest(r, "id"); !errors.Is(err, errs.ErrInvalidFieldValue) {
		t.Fatalf("expected ErrInvalidFieldValue, got %v", err)
	}
}

func TestSwaggerUI(t *testing.T) {
	h := &Handler{}

	r := httptest.NewRequest(http.MethodGet, "/swagger", nil)
	w := httptest.NewRecorder()
	h.swaggerUI(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "SwaggerUIBundle") {
		t.Fatal("expected swagger html page")
	}

	oldWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("failed to get cwd: %v", err)
	}
	if err := os.Chdir("../../.."); err != nil {
		t.Fatalf("failed to chdir to project root: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWd)
	})

	r = httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	w = httptest.NewRecorder()
	h.swaggerUI(w, r)
	if w.Code != http.StatusOK {
		t.Fatalf("expected doc json status 200, got %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "AnalogYouTube API") {
		t.Fatal("expected swagger json content")
	}
}

func TestGetUserIDFromWebSocketRequest(t *testing.T) {
	configs.AppSettings.AuthParams.JwtSecret = "test-secret"
	h := &Handler{}

	token, err := pkg.GenerateToken(22, 15, domain.UserRole, false)
	if err != nil {
		t.Fatalf("failed to generate token: %v", err)
	}

	r := httptest.NewRequest(http.MethodGet, "/ws/chats/1?access_token="+token, nil)
	userID, err := h.getUserIDFromWebSocketRequest(r)
	if err != nil {
		t.Fatalf("getUserIDFromWebSocketRequest returned error: %v", err)
	}
	if userID != 22 {
		t.Fatalf("expected user id 22, got %d", userID)
	}

	r = httptest.NewRequest(http.MethodGet, "/ws/chats/1", nil)
	r.Header.Set(authorizationHeader, "Bearer "+token)
	userID, err = h.getUserIDFromWebSocketRequest(r)
	if err != nil {
		t.Fatalf("header token should work: %v", err)
	}
	if userID != 22 {
		t.Fatalf("expected user id 22, got %d", userID)
	}
}

func TestChatHubStore(t *testing.T) {
	store := &chatHubStore{items: make(map[int]*chatHub)}
	hub := store.getOrCreate(1)
	if hub == nil || hub.chatID != 1 {
		t.Fatalf("unexpected hub: %+v", hub)
	}
	if got := store.get(1); got != hub {
		t.Fatal("expected to get same hub")
	}

	client := &chatClient{userID: 7}
	hub.Register(client)
	if clients := hub.snapshotClients(); len(clients) != 1 || clients[0] != client {
		t.Fatalf("unexpected clients snapshot: %+v", clients)
	}

	store.unregisterClient(1, client)
	if got := store.get(1); got != nil {
		t.Fatal("expected empty hub to be removed")
	}

	store.unregisterClient(404, client)
	broadcastToChat(404, domain.ChatMessage{Text: "ignored"})
}
