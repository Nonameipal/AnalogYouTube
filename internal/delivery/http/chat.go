package httpdelivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"sync"

	"github.com/Nonameipal/AnalogYouTube/internal/delivery/http/dto"
	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/Nonameipal/AnalogYouTube/internal/errs"
	"github.com/Nonameipal/AnalogYouTube/pkg"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type chatHub struct {
	clients map[*websocket.Conn]int
	mu      sync.Mutex
}

var chatHubs = struct {
	items map[int]*chatHub
	mu    sync.Mutex
}{
	items: make(map[int]*chatHub),
}

func getChatHub(chatID int) *chatHub {
	chatHubs.mu.Lock()
	defer chatHubs.mu.Unlock()

	hub, ok := chatHubs.items[chatID]
	if !ok {
		hub = &chatHub{clients: make(map[*websocket.Conn]int)}
		chatHubs.items[chatID] = hub
	}

	return hub
}

func (ctrl *Controller) CreateOrGetChat(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	var input dto.CreateChatRequest
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	chat, err := ctrl.service.CreateOrGetChat(userID, input.UserID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, chat)
}

func (ctrl *Controller) GetMyChats(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	chats, err := ctrl.service.GetUserChats(userID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, chats)
}

func (ctrl *Controller) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	chatID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	messages, err := ctrl.service.GetChatMessages(userID, chatID)
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, messages)
}

func (ctrl *Controller) CreateChatMessage(w http.ResponseWriter, r *http.Request) {
	userID, ok := getUserIDFromContext(r)
	if !ok {
		ctrl.handleError(w, errs.ErrInvalidToken)
		return
	}

	chatID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	var input dto.CreateChatMessageRequest
	if err = json.NewDecoder(r.Body).Decode(&input); err != nil {
		ctrl.handleError(w, errors.Join(errs.ErrInvalidRequestBody, err))
		return
	}

	message, err := ctrl.service.CreateChatMessage(userID, chatID, domain.ChatMessage{Text: input.Text})
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	writeJSON(w, http.StatusCreated, message)
}

func (ctrl *Controller) ChatWebSocket(w http.ResponseWriter, r *http.Request) {
	userID, err := ctrl.getUserIDFromWebSocketRequest(r)
	if err != nil {
		writeJSON(w, http.StatusUnauthorized, CommonError{Error: err.Error()})
		return
	}

	chatID, err := getIDFromRequest(r, "id")
	if err != nil {
		ctrl.handleError(w, err)
		return
	}

	if err = ctrl.service.EnsureUserCanAccessChat(userID, chatID); err != nil {
		ctrl.handleError(w, err)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	defer conn.Close()

	hub := getChatHub(chatID)
	hub.mu.Lock()
	hub.clients[conn] = userID
	hub.mu.Unlock()

	defer func() {
		hub.mu.Lock()
		delete(hub.clients, conn)
		hub.mu.Unlock()
	}()

	for {
		var input dto.WebSocketChatMessageRequest
		if err = conn.ReadJSON(&input); err != nil {
			break
		}

		message, err := ctrl.service.CreateChatMessage(userID, chatID, domain.ChatMessage{Text: input.Text})
		if err != nil {
			_ = conn.WriteJSON(CommonError{Error: err.Error()})
			continue
		}

		broadcastChatMessage(hub, message)
	}
}

func broadcastChatMessage(hub *chatHub, message domain.ChatMessage) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	for conn := range hub.clients {
		if err := conn.WriteJSON(message); err != nil {
			_ = conn.Close()
			delete(hub.clients, conn)
		}
	}
}

func (ctrl *Controller) getUserIDFromWebSocketRequest(r *http.Request) (int, error) {
	tokenString := r.URL.Query().Get("token")
	if tokenString == "" {
		var err error
		tokenString, err = ctrl.extractTokenFromHeader(r, authorizationHeader)
		if err != nil {
			return 0, err
		}
	}

	userID, isRefresh, _, err := pkg.ParseToken(tokenString)
	if err != nil {
		return 0, err
	}
	if isRefresh {
		return 0, errs.ErrInvalidToken
	}

	return userID, nil
}
