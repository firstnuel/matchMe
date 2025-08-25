package websocket

import (
	"context"
	"log"
	"sync"

	"github.com/google/uuid"
)

// =================================================================================
// CHAT HUB IMPLEMENTATION
// Manages WebSocket connections for sending and receiving chat messages.
// =================================================================================

// ChatConnectionGroup represents a group of chat clients in a specific chat room.
type ChatConnectionGroup struct {
	clients map[*Client]bool
	mu      sync.RWMutex
}

func NewChatConnectionGroup() *ChatConnectionGroup {
	return &ChatConnectionGroup{clients: make(map[*Client]bool)}
}

func (cg *ChatConnectionGroup) AddClient(client *Client) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.clients[client] = true
}

func (cg *ChatConnectionGroup) RemoveClient(client *Client) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	delete(cg.clients, client)
}

func (cg *ChatConnectionGroup) GetClientCount() int {
	cg.mu.RLock()
	defer cg.mu.RUnlock()
	return len(cg.clients)
}

// ChatHub maintains the set of active chat clients and broadcasts chat messages.
type ChatHub struct {
	clients     map[*Client]bool
	connections map[uuid.UUID]*ChatConnectionGroup
	register    chan *Client
	unregister  chan *Client
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewChatHub() *ChatHub {
	ctx, cancel := context.WithCancel(context.Background())
	return &ChatHub{
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		connections: make(map[uuid.UUID]*ChatConnectionGroup),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (h *ChatHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			connectionID := client.GetConnectionID()
			if connectionID != nil {
				group, ok := h.connections[*connectionID]
				if !ok {
					group = NewChatConnectionGroup()
					h.connections[*connectionID] = group
				}
				group.AddClient(client)
			}
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				connectionID := client.GetConnectionID()
				if connectionID != nil {
					if group, ok := h.connections[*connectionID]; ok {
						group.RemoveClient(client)
						if group.GetClientCount() == 0 {
							delete(h.connections, *connectionID)
						}
					}
				}
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()

		case <-h.ctx.Done():
			h.mu.Lock()
			for client := range h.clients {
				client.Close()
			}
			h.mu.Unlock()
			return
		}
	}
}

// BroadcastMessage sends a chat message to all clients in a connection except the sender.
func (h *ChatHub) BroadcastMessage(connectionID uuid.UUID, messageEvent MessageEvent, senderUserID uuid.UUID) {
	h.mu.RLock()
	group, ok := h.connections[connectionID]
	h.mu.RUnlock()

	if ok {
		group.mu.RLock()
		defer group.mu.RUnlock()
		for client := range group.clients {
			if client.userID != senderUserID {
				client.SendMessage(EventMessageNew, messageEvent)
			}
		}
	}
}

// GetConnectionUsers returns a slice of user IDs for a given connection.
func (h *ChatHub) GetConnectionUsers(connectionID uuid.UUID) ([]uuid.UUID, bool) {
	h.mu.RLock()
	group, ok := h.connections[connectionID]
	h.mu.RUnlock()

	if !ok {
		return nil, false
	}

	group.mu.RLock()
	defer group.mu.RUnlock()

	users := make([]uuid.UUID, 0, len(group.clients))
	for client := range group.clients {
		users = append(users, client.userID)
	}

	return users, true
}

// BroadcastEvent sends any event to all clients in a connection except the sender.
func (h *ChatHub) BroadcastEvent(connectionID uuid.UUID, eventType EventType, data interface{}, senderUserID uuid.UUID) {
	h.mu.RLock()
	group, ok := h.connections[connectionID]
	h.mu.RUnlock()

	if ok {
		group.mu.RLock()
		defer group.mu.RUnlock()
		for client := range group.clients {
			if client.userID != senderUserID {
				client.SendMessage(eventType, data)
			}
		}
	}
}

func (h *ChatHub) Shutdown() {
	h.cancel()
}

// =================================================================================
// TYPING HUB IMPLEMENTATION
// Manages WebSocket connections for sending and receiving typing indicators.
// =================================================================================

// TypingConnectionGroup represents a group of typing clients in a specific chat room.
type TypingConnectionGroup struct {
	clients map[*Client]bool
	mu      sync.RWMutex
}

func NewTypingConnectionGroup() *TypingConnectionGroup {
	return &TypingConnectionGroup{clients: make(map[*Client]bool)}
}

func (cg *TypingConnectionGroup) AddClient(client *Client) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	cg.clients[client] = true
}

func (cg *TypingConnectionGroup) RemoveClient(client *Client) {
	cg.mu.Lock()
	defer cg.mu.Unlock()
	delete(cg.clients, client)
}

func (cg *TypingConnectionGroup) GetClientCount() int {
	cg.mu.RLock()
	defer cg.mu.RUnlock()
	return len(cg.clients)
}

// TypingHub maintains the set of active typing clients and broadcasts typing events.
type TypingHub struct {
	clients     map[*Client]bool
	connections map[uuid.UUID]*TypingConnectionGroup
	register    chan *Client
	unregister  chan *Client
	mu          sync.RWMutex
	ctx         context.Context
	cancel      context.CancelFunc
}

func NewTypingHub() *TypingHub {
	ctx, cancel := context.WithCancel(context.Background())
	return &TypingHub{
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		clients:     make(map[*Client]bool),
		connections: make(map[uuid.UUID]*TypingConnectionGroup),
		ctx:         ctx,
		cancel:      cancel,
	}
}

func (h *TypingHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			connectionID := client.GetConnectionID()
			if connectionID != nil {
				group, ok := h.connections[*connectionID]
				if !ok {
					group = NewTypingConnectionGroup()
					h.connections[*connectionID] = group
				}
				group.AddClient(client)
			}
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				connectionID := client.GetConnectionID()
				if connectionID != nil {
					if group, ok := h.connections[*connectionID]; ok {
						group.RemoveClient(client)
						if group.GetClientCount() == 0 {
							delete(h.connections, *connectionID)
						}
					}
				}
				delete(h.clients, client)
				client.Close()
			}
			h.mu.Unlock()

		case <-h.ctx.Done():
			h.mu.Lock()
			for client := range h.clients {
				client.Close()
			}
			h.mu.Unlock()
			return
		}
	}
}

// BroadcastTypingIndicator sends a typing event to all clients in a connection except the sender.
func (h *TypingHub) BroadcastTypingIndicator(connectionID uuid.UUID, typingEvent TypingEvent, senderUserID uuid.UUID) {
	h.mu.RLock()
	group, ok := h.connections[connectionID]
	h.mu.RUnlock()

	if ok {
		group.mu.RLock()
		defer group.mu.RUnlock()
		for client := range group.clients {
			if client.userID != senderUserID {
				client.SendMessage(EventMessageTyping, typingEvent)
			}
		}
	}
}

func (h *TypingHub) Shutdown() {
	h.cancel()
}

// =================================================================================
// STATUS HUB IMPLEMENTATION
// Manages WebSocket connections for user presence and direct notifications.
// (e.g., connection requests, user online/offline status)
// =================================================================================

// StatusHub maintains the set of active clients for status and direct messaging.
type StatusHub struct {
	// All connected clients.
	clients map[*Client]bool

	// Maps a userID to a set of their active clients.
	// A user might have multiple status connections (e.g., from different browser tabs).
	clientsByUser map[uuid.UUID]map[*Client]bool

	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
}

// NewStatusHub creates a new StatusHub.
func NewStatusHub() *StatusHub {
	ctx, cancel := context.WithCancel(context.Background())
	return &StatusHub{
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		clients:       make(map[*Client]bool),
		clientsByUser: make(map[uuid.UUID]map[*Client]bool),
		ctx:           ctx,
		cancel:        cancel,
	}
}

// Run starts the hub's event loop.
func (h *StatusHub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client] = true
			if h.clientsByUser[client.userID] == nil {
				h.clientsByUser[client.userID] = make(map[*Client]bool)
			}
			h.clientsByUser[client.userID][client] = true
			log.Printf("✅ Status client registered for user %s", client.userID)
			h.mu.Unlock()

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				// Remove the client from the user-specific map.
				if userClients, ok := h.clientsByUser[client.userID]; ok {
					delete(userClients, client)
					// If the user has no more active status connections, remove their entry.
					if len(userClients) == 0 {
						delete(h.clientsByUser, client.userID)
					}
				}
				client.Close()
				log.Printf("🔌 Status client unregistered for user %s", client.userID)
			}
			h.mu.Unlock()

		case <-h.ctx.Done():
			h.mu.Lock()
			log.Println("Shutting down status hub...")
			for client := range h.clients {
				client.Close()
			}
			h.mu.Unlock()
			return
		}
	}
}

// BroadcastToUser sends a direct message to all connections for a specific user.
func (h *StatusHub) BroadcastToUser(userID uuid.UUID, eventType EventType, data interface{}) {
	h.mu.RLock()
	// Find the clients for the target user and create a copy to avoid holding the lock during send.
	var clientsToSend []*Client
	if userClients, ok := h.clientsByUser[userID]; ok {
		clientsToSend = make([]*Client, 0, len(userClients))
		for c := range userClients {
			clientsToSend = append(clientsToSend, c)
		}
	}
	h.mu.RUnlock()

	// Send the message to the copied list of clients.
	for _, client := range clientsToSend {
		client.SendMessage(eventType, data)
	}
}

// broadcastUserStatus broadcasts a user's status change to all their connections.
func (h *StatusHub) broadcastUserStatus(userID uuid.UUID, status string) {
	h.mu.RLock()
	// Create a copy of the connection IDs to avoid holding the lock during broadcast.
	// connectionIDs := make([]uuid.UUID, 0, len(h.userConnections[userID]))
	// for connID := range h.userConnections[userID] {
	// 	connectionIDs = append(connectionIDs, connID)
	// }
	// h.mu.RUnlock()

	// statusEvent := UserStatusEvent{
	// 	UserID:       userID.String(),
	// 	Status:       status,
	// 	LastActivity: time.Now().UTC().Format(time.RFC3339),
	// }

	// eventType := EventUserOnline // Default event type
	// if status == "offline" {
	// 	eventType = EventUserOffline
	// }

	// for _, connID := range connectionIDs {
	// Broadcast to everyone in that connection, excluding the user whose status changed.
	// NOTE: This assumes you have a ChatHub instance available or a way to access its broadcast method.
	// A cleaner way would be to pass the ChatHub to this method or make it available to the StatusHub.
	// For now, this is a placeholder for the broadcast logic.
	// log.Printf("Need to broadcast status for user %s to connection %s", userID, connID)
	// }
}

// GetOnlineUsers returns a slice of unique user IDs of all clients connected to this hub.
func (h *StatusHub) GetOnlineUsers() []uuid.UUID {
	h.mu.RLock()
	defer h.mu.RUnlock()

	userIDs := make([]uuid.UUID, 0, len(h.clientsByUser))
	for userID := range h.clientsByUser {
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

// IsUserOnline checks if a user has at least one active status connection.
func (h *StatusHub) IsUserOnline(userID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	_, ok := h.clientsByUser[userID]
	return ok
}

// Shutdown gracefully stops the hub and closes all connections.
func (h *StatusHub) Shutdown() {
	h.cancel()
}
