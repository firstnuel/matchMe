package websocket

import (
	"log"
	"net/http"

	"match-me/api/middleware"
	"match-me/config"
	"match-me/internal/repositories/connections"
	userUsecase "match-me/internal/usecases/user"
	wscore "match-me/internal/websocket"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// WebSocketHandler handles WebSocket connections
type WebSocketHandler struct {
	hub            *wscore.Hub
	connectionRepo connections.ConnectionRepository
	UserUsecase    userUsecase.UserUsecase
	cfg            *config.Config
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(hub *wscore.Hub,
	connectionRepo connections.ConnectionRepository,
	userUsecase userUsecase.UserUsecase,
	cfg *config.Config,
) *WebSocketHandler {
	return &WebSocketHandler{
		hub:            hub,
		UserUsecase:    userUsecase,
		connectionRepo: connectionRepo,
		cfg:            cfg,
	}
}

// HandleChatConnection handles WebSocket connections for chat
func (h *WebSocketHandler) HandleChatConnection(c *gin.Context) {
	// Get user ID from middleware
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get connection ID from URL parameter
	connectionIDStr := c.Param("connectionId")
	connectionID, err := uuid.Parse(connectionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}

	// Verify that the user is part of this connection
	if !h.validateConnectionAccess(c, user.ID, connectionID) {
		return
	}

	// Upgrade to WebSocket and create client
	wscore.ServeWS(h.hub, c, user.ID)

	// Add client to connection group
	h.hub.AddClientToConnection(user.ID, connectionID)

	log.Printf("Chat WebSocket connection established for user %s in connection %s", user.ID, connectionID)
}

// HandleStatusConnection handles WebSocket connections for user status
func (h *WebSocketHandler) HandleStatusConnection(c *gin.Context) {
	// Get user ID from middleware
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Upgrade to WebSocket and create client
	wscore.ServeWS(h.hub, c, user.ID)

	log.Printf("Status WebSocket connection established for user %s", user.ID)
}

// HandleTypingConnection handles WebSocket connections for typing indicators
func (h *WebSocketHandler) HandleTypingConnection(c *gin.Context) {
	// Get user ID from middleware
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Get connection ID from URL parameter
	connectionIDStr := c.Param("connectionId")
	connectionID, err := uuid.Parse(connectionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}

	// Verify that the user is part of this connection
	if !h.validateConnectionAccess(c, user.ID, connectionID) {
		return
	}

	// Upgrade to WebSocket and create client
	wscore.ServeWS(h.hub, c, user.ID)

	// Add client to connection group for typing
	h.hub.AddClientToConnection(user.ID, connectionID)

	log.Printf("Typing WebSocket connection established for user %s in connection %s", user.ID, connectionID)
}

// validateConnectionAccess verifies that a user has access to a connection
func (h *WebSocketHandler) validateConnectionAccess(c *gin.Context, userID, connectionID uuid.UUID) bool {
	// Get the connection from repository
	connection, err := h.connectionRepo.GetConnection(c.Request.Context(), connectionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Connection not found"})
		return false
	}

	// Verify the connection is active
	if connection.Status != "connected" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Connection is not active"})
		return false
	}

	// Verify user is part of this connection
	if connection.UserAID != userID && connection.UserBID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this connection"})
		return false
	}

	return true
}

// GetOnlineUsers returns list of online users (for debugging/admin)
func (h *WebSocketHandler) GetOnlineUsers(c *gin.Context) {
	onlineUsers := h.hub.GetOnlineUsers()
	c.JSON(http.StatusOK, gin.H{
		"online_users": onlineUsers,
		"count":        len(onlineUsers),
	})
}

// GetConnectionStatus returns status of users in a connection
func (h *WebSocketHandler) GetConnectionStatus(c *gin.Context) {
	// Get connection ID from URL parameter
	connectionIDStr := c.Param("connectionId")
	connectionID, err := uuid.Parse(connectionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}

	// Get user ID from middleware
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	// Verify access to connection
	if !h.validateConnectionAccess(c, user.ID, connectionID) {
		return
	}

	// Get connection users
	connectionUsers := h.hub.GetConnectionUsers(connectionID)

	// Build status map
	userStatuses := make(map[string]bool)
	for _, uid := range connectionUsers {
		userStatuses[uid.String()] = h.hub.IsUserOnline(uid)
	}

	c.JSON(http.StatusOK, gin.H{
		"connection_id": connectionID,
		"users":         userStatuses,
	})
}
