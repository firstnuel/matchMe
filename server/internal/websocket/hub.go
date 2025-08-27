package websocket

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
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
	// Start stale connection cleanup goroutine
	go h.cleanupStaleConnections()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()

			// Initialize user's client map if it doesn't exist
			if h.clientsByUser[client.userID] == nil {
				h.clientsByUser[client.userID] = make(map[*Client]bool)
			}

			// Check if user was offline BEFORE adding the new client
			wasOffline := len(h.clientsByUser[client.userID]) == 0

			// Add client to global clients map
			h.clients[client] = true

			// Add client to user's client map
			h.clientsByUser[client.userID][client] = true
			log.Printf("✅ Status client registered for user %s (total connections: %d)", client.userID, len(h.clientsByUser[client.userID]))
			h.mu.Unlock()

			// 1. Send the list of already-online users directly to the new client.
			h.SendInitialUserStatuses(client)

			// 2. Broadcast the new user's "online" status to everyone else.
			if wasOffline {
				h.BroadcastUserStatus(client.userID, "online")
			}

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				userID := client.userID
				// Remove the client from the user-specific map.
				if userClients, ok := h.clientsByUser[userID]; ok {
					delete(userClients, client)
					// Check if user will be offline after removing this client
					willBeOffline := len(userClients) == 0
					// If the user has no more active status connections, remove their entry.
					if willBeOffline {
						delete(h.clientsByUser, userID)
					}
					client.Close()
					log.Printf("🔌 Status client unregistered for user %s", userID)
					h.mu.Unlock()

					// Broadcast user offline status if they have no more connections
					if willBeOffline {
						h.BroadcastUserStatus(userID, "offline")
					}
				} else {
					h.mu.Unlock()
				}
			} else {
				h.mu.Unlock()
			}

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

// BroadcastUserStatus broadcasts a user's status change to all connected status clients except the user themselves.
func (h *StatusHub) BroadcastUserStatus(userID uuid.UUID, status string) {
	h.mu.RLock()
	// Create a copy of all clients except the user whose status is changing
	clientsToSend := make([]*Client, 0, len(h.clients))
	for client := range h.clients {
		// Exclude the user whose status is changing to avoid feedback loops
		if client.userID != userID {
			clientsToSend = append(clientsToSend, client)
		}
	}
	h.mu.RUnlock()

	// Don't broadcast if no other clients are connected
	if len(clientsToSend) == 0 {
		log.Printf("📡 No other clients to broadcast %s status for user %s", status, userID)
		return
	}

	// Create the status event
	statusEvent := UserStatusEvent{
		UserID:       userID,
		Status:       status,
		LastActivity: time.Now().UTC(),
	}

	// Determine the correct event type
	var eventType EventType
	switch status {
	case "online":
		eventType = EventUserOnline
	case "offline":
		eventType = EventUserOffline
	case "away":
		eventType = EventUserAway
	default:
		eventType = EventUserOnline
	}

	// Broadcast to all status clients except the user themselves
	for _, client := range clientsToSend {
		client.SendMessage(eventType, statusEvent)
	}

	log.Printf("📡 Broadcasted %s status for user %s to %d other clients", status, userID, len(clientsToSend))
}

// SetUserAway marks a user as away and broadcasts the status change
func (h *StatusHub) SetUserAway(userID uuid.UUID) {
	h.mu.RLock()
	_, userExists := h.clientsByUser[userID]
	h.mu.RUnlock()

	// Only broadcast if user is currently online
	if userExists {
		h.BroadcastUserStatus(userID, "away")
	}
}

// SetUserOnline marks a user as online and broadcasts the status change
func (h *StatusHub) SetUserOnline(userID uuid.UUID) {
	h.mu.RLock()
	_, userExists := h.clientsByUser[userID]
	h.mu.RUnlock()

	// Only broadcast if user is currently online
	if userExists {
		h.BroadcastUserStatus(userID, "online")
	}
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

// SendInitialUserStatuses sends the current status of all online users to a single new client.
func (h *StatusHub) SendInitialUserStatuses(newClient *Client) {
	h.mu.RLock()

	// Create a list of status events for all users who are already online.
	// We exclude the new client themselves, as they know they are online.
	onlineUsers := make([]UserStatusEvent, 0, len(h.clientsByUser))
	for userID := range h.clientsByUser {
		if userID != newClient.userID {
			onlineUsers = append(onlineUsers, UserStatusEvent{
				UserID:       userID,
				Status:       "online",
				LastActivity: time.Now().UTC(), // Or fetch last known activity
			})
		}
	}
	h.mu.RUnlock()

	// Send this snapshot of online users to the new client.
	if len(onlineUsers) > 0 {
		newClient.SendMessage(EventUserStatusInitial, onlineUsers)
	}
}

// cleanupStaleConnections periodically removes inactive connections.
func (h *StatusHub) cleanupStaleConnections() {
	ticker := time.NewTicker(60 * time.Second) // Increased from 30s to 60s for less aggressive cleanup
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			h.mu.Lock()
			var staleClients []*Client
			for client := range h.clients {
				// More conservative stale detection - only remove truly dead connections
				if client.IsStale() {
					// Double-check by testing the WebSocket connection
					if err := client.conn.WriteControl(websocket.PingMessage, []byte{}, time.Now().Add(5*time.Second)); err != nil {
						staleClients = append(staleClients, client)
					} else {
						// Connection is still alive, update activity
						client.UpdateActivity()
					}
				}
			}
			h.mu.Unlock()

			// Only unregister clients that are truly stale
			for _, client := range staleClients {
				select {
				case h.unregister <- client:
					log.Printf("🧹 Cleaned up truly stale connection for user %s", client.userID)
				default:
					// Channel is full, skip this cleanup cycle
				}
			}

		case <-h.ctx.Done():
			return
		}
	}
}

// Shutdown gracefully stops the hub and closes all connections.
func (h *StatusHub) Shutdown() {
	h.cancel()
}
