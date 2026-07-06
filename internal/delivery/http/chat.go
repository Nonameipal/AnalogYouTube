package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/Nonameipal/AnalogYouTube/internal/delivery/http/dto"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/pkg"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 10 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *Handler) SendChatRequest(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	var input dto.CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	request, err := h.service.SendChatRequest(userID, input.UserID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, request)
}

func (h *Handler) GetMyChats(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	chats, err := h.service.GetUserChats(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, chats)
}

func (h *Handler) GetIncomingChatRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	requests, err := h.service.GetIncomingChatRequests(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, requests)
}

func (h *Handler) GetOutgoingChatRequests(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	requests, err := h.service.GetOutgoingChatRequests(userID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, requests)
}

func (h *Handler) AcceptChatRequest(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	requestID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	result, err := h.service.AcceptChatRequest(userID, requestID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, result)
}

func (h *Handler) RejectChatRequest(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	requestID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	request, err := h.service.RejectChatRequest(userID, requestID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, request)
}

func (h *Handler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	chatID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	messages, err := h.service.GetChatMessages(userID, chatID)
	if err != nil {
		h.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, messages)
}

func (h *Handler) CreateChatMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		h.handleError(w, errs.ErrInvalidToken)
		return
	}

	chatID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	var input dto.CreateChatMessageRequest
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		h.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	message, err := h.service.CreateChatMessage(userID, chatID, domain.ChatMessage{Text: input.MessageText()})
	if err != nil {
		h.handleError(w, err)
		return
	}

	broadcastToChat(chatID, message)

	writeJSON(w, http.StatusCreated, message)
}

func (h *Handler) ChatWebSocket(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromWebSocketRequest(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, CommonError{Error: err.Error()})
		return
	}

	chatID, err := getIDFromRequest(r, "id")
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err = h.service.EnsureUserCanAccessChat(userID, chatID); err != nil {
		h.handleError(w, err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	conn.SetReadLimit(4096)

	client := newChatClient(userID, conn)
	chatRooms.getOrCreate(chatID).Register(client)
	defer chatRooms.unregisterClient(chatID, client)

	for {
		var input dto.WebSocketChatMessageRequest
		if err = conn.ReadJSON(&input); err != nil {
			break
		}

		message, err := h.service.CreateChatMessage(userID, chatID, domain.ChatMessage{Text: input.MessageText()})
		if err != nil {
			_ = client.WriteJSON(CommonError{Error: err.Error()})
			continue
		}

		broadcastToChat(chatID, message)
	}
}

func (h *Handler) getUserIDFromWebSocketRequest(r *http.Request) (int, error) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		tokenString = r.URL.Query().Get("access_token")
	}
	if tokenString == "" {
		var err error
		tokenString, err = h.extractTokenFromHeader(r, authorizationHeader)
		if err != nil {
			return 0, err
		}
	}

	tokenString = strings.TrimSpace(strings.TrimPrefix(tokenString, "Bearer "))

	userID, isRefresh, _, err := pkg.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	if isRefresh {
		return 0, errs.ErrInvalidToken
	}

	return userID, nil
}
