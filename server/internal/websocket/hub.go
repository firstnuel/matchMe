package websocket

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/google/uuid"
)

// Hub maintains the set of active clients and broadcasts messages to the clients.
type Hub struct {
	// Registered clients by user ID
	clients map[uuid.UUID]*Client

	// Clients grouped by connection ID for chat
	connections map[uuid.UUID]map[uuid.UUID]*Client

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Mutex for concurrent access
	mu sync.RWMutex

	// Context for shutdown
	ctx    context.Context
	cancel context.CancelFunc
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	ctx, cancel := context.WithCancel(context.Background())
	return &Hub{
		clients:     make(map[uuid.UUID]*Client),
		connections: make(map[uuid.UUID]map[uuid.UUID]*Client),
		broadcast:   make(chan []byte),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		ctx:         ctx,
		cancel:      cancel,
	}
}

// Run starts the hub and handles client registration/unregistration
func (h *Hub) Run() {
	// Start cleanup routine
	go h.cleanupRoutine()

	for {
		select {
		case <-h.ctx.Done():
			return

		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case message := <-h.broadcast:
			// Handle broadcast messages if needed
			log.Printf("Broadcasting message: %s", message)
		}
	}
}

// registerClient registers a new client
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Close existing client connection if exists
	if existingClient, exists := h.clients[client.userID]; exists {
		existingClient.Close()
	}

	// Register the new client
	h.clients[client.userID] = client

	// Add to connection group if this is a chat client
	if connectionID := client.GetConnectionID(); connectionID != nil {
		if h.connections[*connectionID] == nil {
			h.connections[*connectionID] = make(map[uuid.UUID]*Client)
		}
		h.connections[*connectionID][client.userID] = client
	}

	log.Printf("Client registered: user %s", client.userID)

	// Broadcast user online status
	h.broadcastUserStatus(client.userID, "online")
}

// unregisterClient unregisters a client
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	// Remove from clients map
	delete(h.clients, client.userID)

	// Remove from connection group
	if connectionID := client.GetConnectionID(); connectionID != nil {
		if connectionClients, exists := h.connections[*connectionID]; exists {
			delete(connectionClients, client.userID)
			// Clean up empty connection groups
			if len(connectionClients) == 0 {
				delete(h.connections, *connectionID)
			}
		}
	}

	client.Close()
	log.Printf("Client unregistered: user %s", client.userID)

	// Broadcast user offline status
	h.broadcastUserStatus(client.userID, "offline")
}

// AddClientToConnection adds a client to a specific connection group
func (h *Hub) AddClientToConnection(userID, connectionID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	client, exists := h.clients[userID]
	if !exists {
		return
	}

	client.SetConnectionID(connectionID)

	if h.connections[connectionID] == nil {
		h.connections[connectionID] = make(map[uuid.UUID]*Client)
	}
	h.connections[connectionID][userID] = client
}

// BroadcastToUser sends a message to a specific user
func (h *Hub) BroadcastToUser(userID uuid.UUID, eventType EventType, data interface{}) {
	h.mu.RLock()
	client, exists := h.clients[userID]
	h.mu.RUnlock()

	if exists && client.IsActive() {
		client.SendMessage(eventType, data)
	}
}

// broadcastToConnection sends a message to all users in a connection except the sender
func (h *Hub) broadcastToConnection(connectionID uuid.UUID, eventType EventType, data interface{}, senderID uuid.UUID) {
	h.mu.RLock()
	connectionClients, exists := h.connections[connectionID]
	h.mu.RUnlock()

	if !exists {
		return
	}

	for userID, client := range connectionClients {
		if userID != senderID && client.IsActive() {
			client.SendMessage(eventType, data)
		}
	}
}

// BroadcastMessageToConnection broadcasts a new message to connection participants
func (h *Hub) BroadcastMessageToConnection(connectionID uuid.UUID, messageEvent MessageEvent) {
	h.broadcastToConnection(connectionID, EventMessageNew, messageEvent, messageEvent.SenderID)
}

// BroadcastMessageRead broadcasts message read status to sender
func (h *Hub) BroadcastMessageRead(connectionID uuid.UUID, readEvent MessageReadEvent) {
	// Send read receipt to the message sender (not the reader)
	h.mu.RLock()
	connectionClients, exists := h.connections[connectionID]
	h.mu.RUnlock()

	if exists {
		for userID, client := range connectionClients {
			if userID != readEvent.ReadBy && client.IsActive() {
				client.SendMessage(EventMessageRead, readEvent)
			}
		}
	}
}

// broadcastUserStatus broadcasts user status changes to their active connections
func (h *Hub) broadcastUserStatus(userID uuid.UUID, status string) {
	statusEvent := UserStatusEvent{
		UserID:       userID,
		Status:       status,
		LastActivity: time.Now(),
	}

	// Find all connections this user is part of and notify other participants
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, connectionClients := range h.connections {
		if _, isInConnection := connectionClients[userID]; isInConnection {
			// Broadcast to other users in this connection
			for otherUserID, client := range connectionClients {
				if otherUserID != userID && client.IsActive() {
					client.SendMessage(EventUserOnline, statusEvent)
				}
			}
		}
	}
}

// GetOnlineUsers returns list of online user IDs
func (h *Hub) GetOnlineUsers() []uuid.UUID {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]uuid.UUID, 0, len(h.clients))
	for userID, client := range h.clients {
		if client.IsActive() {
			users = append(users, userID)
		}
	}
	return users
}

// IsUserOnline checks if a user is currently online
func (h *Hub) IsUserOnline(userID uuid.UUID) bool {
	h.mu.RLock()
	defer h.mu.RUnlock()

	client, exists := h.clients[userID]
	return exists && client.IsActive()
}

// GetConnectionUsers returns users in a specific connection
func (h *Hub) GetConnectionUsers(connectionID uuid.UUID) []uuid.UUID {
	h.mu.RLock()
	defer h.mu.RUnlock()

	connectionClients, exists := h.connections[connectionID]
	if !exists {
		return nil
	}

	users := make([]uuid.UUID, 0, len(connectionClients))
	for userID := range connectionClients {
		users = append(users, userID)
	}
	return users
}

// cleanupRoutine periodically cleans up inactive clients
func (h *Hub) cleanupRoutine() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-h.ctx.Done():
			return
		case <-ticker.C:
			h.cleanupInactiveClients()
		}
	}
}

// cleanupInactiveClients removes clients that have been inactive for too long
func (h *Hub) cleanupInactiveClients() {
	h.mu.Lock()
	defer h.mu.Unlock()

	inactiveThreshold := time.Now().Add(-10 * time.Minute)

	for userID, client := range h.clients {
		if !client.IsActive() || client.lastActivity.Before(inactiveThreshold) {
			log.Printf("Cleaning up inactive client: user %s", userID)
			h.unregister <- client
		}
	}
}

// Shutdown gracefully shuts down the hub
func (h *Hub) Shutdown() {
	h.cancel()

	// Close all client connections
	h.mu.Lock()
	for _, client := range h.clients {
		client.Close()
	}
	h.mu.Unlock()
}
