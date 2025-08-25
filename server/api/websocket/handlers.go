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
	chatHub        *wscore.ChatHub
	typingHub      *wscore.TypingHub
	statusHub      *wscore.StatusHub
	connectionRepo connections.ConnectionRepository
	UserUsecase    userUsecase.UserUsecase
	cfg            *config.Config
}

// NewWebSocketHandler creates a new WebSocket handler
func NewWebSocketHandler(chatHub *wscore.ChatHub, typingHub *wscore.TypingHub, statusHub *wscore.StatusHub,
	connectionRepo connections.ConnectionRepository,
	userUsecase userUsecase.UserUsecase,
	cfg *config.Config,
) *WebSocketHandler {
	return &WebSocketHandler{
		chatHub:        chatHub,
		typingHub:      typingHub,
		statusHub:      statusHub,
		UserUsecase:    userUsecase,
		connectionRepo: connectionRepo,
		cfg:            cfg,
	}
}

// HandleChatConnection handles WebSocket connections for chat
func (h *WebSocketHandler) HandleChatConnection(c *gin.Context) {
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	connectionIDStr := c.Param("connectionId")
	connectionID, err := uuid.Parse(connectionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}

	if !h.validateConnectionAccess(c, user.ID, connectionID) {
		return
	}

	log.Printf("🔌 Starting Chat WebSocket upgrade for user %s in connection %s", user.ID, connectionID)
	// FIX: Call the specific function for the ChatHub.
	wscore.ServeChatWS(h.chatHub, c, user.ID, connectionID)
}

// HandleStatusConnection handles WebSocket connections for user status
func (h *WebSocketHandler) HandleStatusConnection(c *gin.Context) {
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	log.Printf("📊 Starting Status WebSocket upgrade for user %s", user.ID)
	// FIX: Call the specific function for the StatusHub.
	wscore.ServeStatusWS(h.statusHub, c, user.ID)
}

// HandleTypingConnection handles WebSocket connections for typing indicators
func (h *WebSocketHandler) HandleTypingConnection(c *gin.Context) {
	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	connectionIDStr := c.Param("connectionId")
	connectionID, err := uuid.Parse(connectionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}

	if !h.validateConnectionAccess(c, user.ID, connectionID) {
		return
	}

	log.Printf("⌨️ Starting Typing WebSocket upgrade for user %s in connection %s", user.ID, connectionID)
	// FIX: Call the specific function for the TypingHub.
	wscore.ServeTypingWS(h.typingHub, c, user.ID, connectionID)
}

// validateConnectionAccess verifies that a user has access to a connection
func (h *WebSocketHandler) validateConnectionAccess(c *gin.Context, userID, connectionID uuid.UUID) bool {
	connection, err := h.connectionRepo.GetConnection(c.Request.Context(), connectionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Connection not found"})
		return false
	}

	if connection.Status != "connected" {
		c.JSON(http.StatusForbidden, gin.H{"error": "Connection is not active"})
		return false
	}

	if connection.UserAID != userID && connection.UserBID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "Access denied to this connection"})
		return false
	}

	return true
}

// GetOnlineUsers returns list of online users (for debugging/admin)
func (h *WebSocketHandler) GetOnlineUsers(c *gin.Context) {
	onlineUsers := h.statusHub.GetOnlineUsers()
	c.JSON(http.StatusOK, gin.H{
		"online_users": onlineUsers,
		"count":        len(onlineUsers),
	})
}

// GetConnectionStatus returns status of users in a connection
func (h *WebSocketHandler) GetConnectionStatus(c *gin.Context) {
	connectionIDStr := c.Param("connectionId")
	connectionID, err := uuid.Parse(connectionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid connection ID"})
		return
	}

	user, exists := middleware.GetUserFromGinContext(c)
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	if !h.validateConnectionAccess(c, user.ID, connectionID) {
		return
	}

	connection, err := h.connectionRepo.GetConnection(c.Request.Context(), connectionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Connection not found"})
		return
	}

	userStatuses := make(map[string]bool)
	userStatuses[connection.UserAID.String()] = h.statusHub.IsUserOnline(connection.UserAID)
	userStatuses[connection.UserBID.String()] = h.statusHub.IsUserOnline(connection.UserBID)

	c.JSON(http.StatusOK, gin.H{
		"connection_id": connectionID,
		"users":         userStatuses,
	})
}
