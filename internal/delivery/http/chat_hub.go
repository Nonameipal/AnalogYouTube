package httpdelivery

import (
	"sync"

	"github.com/Nonameipal/AnalogYouTube/internal/domain"
	"github.com/gorilla/websocket"
)

type chatClient struct {
	userID int
	conn   *websocket.Conn
	mu     sync.Mutex
}

func newChatClient(userID int, conn *websocket.Conn) *chatClient {
	return &chatClient{
		userID: userID,
		conn:   conn,
	}
}

func (c *chatClient) WriteJSON(value any) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	return c.conn.WriteJSON(value)
}

func (c *chatClient) Close() error {
	return c.conn.Close()
}

type chatHub struct {
	chatID  int
	clients map[*chatClient]struct{}
	mu      sync.RWMutex
}

func newChatHub(chatID int) *chatHub {
	return &chatHub{
		chatID:  chatID,
		clients: make(map[*chatClient]struct{}),
	}
}

func (h *chatHub) Register(client *chatClient) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[client] = struct{}{}
}

func (h *chatHub) Unregister(client *chatClient) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	delete(h.clients, client)
	return len(h.clients) == 0
}

func (h *chatHub) Broadcast(message domain.ChatMessage) {
	clients := h.snapshotClients()

	for _, client := range clients {
		if err := client.WriteJSON(message); err != nil {
			_ = client.Close()
			chatRooms.unregisterClient(h.chatID, client)
		}
	}
}

func (h *chatHub) snapshotClients() []*chatClient {
	h.mu.RLock()
	defer h.mu.RUnlock()

	clients := make([]*chatClient, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}

	return clients
}

type chatHubStore struct {
	items map[int]*chatHub
	mu    sync.Mutex
}

var chatRooms = &chatHubStore{
	items: make(map[int]*chatHub),
}

func (s *chatHubStore) getOrCreate(chatID int) *chatHub {
	s.mu.Lock()
	defer s.mu.Unlock()

	hub, ok := s.items[chatID]
	if !ok {
		hub = newChatHub(chatID)
		s.items[chatID] = hub
	}

	return hub
}

func (s *chatHubStore) get(chatID int) *chatHub {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.items[chatID]
}

func (s *chatHubStore) unregisterClient(chatID int, client *chatClient) {
	s.mu.Lock()
	hub := s.items[chatID]
	s.mu.Unlock()

	if hub == nil {
		return
	}

	if hub.Unregister(client) {
		s.mu.Lock()
		if s.items[chatID] == hub {
			delete(s.items, chatID)
		}
		s.mu.Unlock()
	}
}

func broadcastToChat(chatID int, message domain.ChatMessage) {
	hub := chatRooms.get(chatID)
	if hub == nil {
		return
	}

	hub.Broadcast(message)
}
