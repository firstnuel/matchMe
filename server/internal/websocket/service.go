package websocket

import (
	"match-me/internal/models"
	"time"

	"github.com/google/uuid"
)

// WebSocketService provides methods to integrate WebSocket with business logic
type WebSocketService struct {
	hub *Hub
}

// NewWebSocketService creates a new WebSocket service
func NewWebSocketService(hub *Hub) *WebSocketService {
	return &WebSocketService{
		hub: hub,
	}
}

// BroadcastNewMessage broadcasts a new message to connection participants
func (s *WebSocketService) BroadcastNewMessage(message *models.Message) {
	if message == nil {
		return
	}

	messageEvent := MessageEvent{
		Message:      message,
		ConnectionID: message.ConnectionID,
		SenderID:     message.SenderID,
		ReceiverID:   message.ReceiverID,
	}

	s.hub.BroadcastMessageToConnection(message.ConnectionID, messageEvent)
}

// BroadcastMessageRead broadcasts message read status to the sender
func (s *WebSocketService) BroadcastMessageRead(messageID, connectionID, readByUserID uuid.UUID) {
	readEvent := MessageReadEvent{
		MessageID:    messageID,
		ConnectionID: connectionID,
		ReadBy:       readByUserID,
		ReadAt:       time.Now(),
	}

	s.hub.BroadcastMessageRead(connectionID, readEvent)
}

// BroadcastConnectionMessagesRead broadcasts that all messages in a connection were read
func (s *WebSocketService) BroadcastConnectionMessagesRead(connectionID, readByUserID uuid.UUID, messageCount int) {
	// Create a generic read event for the connection
	readEvent := MessageReadEvent{
		MessageID:    uuid.Nil, // No specific message ID for bulk read
		ConnectionID: connectionID,
		ReadBy:       readByUserID,
		ReadAt:       time.Now(),
	}

	s.hub.BroadcastMessageRead(connectionID, readEvent)
}

// BroadcastTypingIndicator broadcasts typing status to connection participants
func (s *WebSocketService) BroadcastTypingIndicator(connectionID, userID uuid.UUID, isTyping bool) {
	typingEvent := TypingEvent{
		ConnectionID: connectionID,
		UserID:       userID,
		IsTyping:     isTyping,
		UpdatedAt:    time.Now(),
	}

	s.hub.broadcastToConnection(connectionID, EventMessageTyping, typingEvent, userID)
}

// BroadcastConnectionRequest broadcasts a new connection request to the receiver
func (s *WebSocketService) BroadcastConnectionRequest(request *models.ConnectionRequest) {
	if request == nil {
		return
	}

	requestEvent := ConnectionRequestEvent{
		Request: request,
		Action:  "new",
	}

	s.hub.BroadcastToUser(request.ReceiverID, EventConnectionRequest, requestEvent)
}

// BroadcastConnectionAccepted broadcasts that a connection request was accepted
func (s *WebSocketService) BroadcastConnectionAccepted(request *models.ConnectionRequest, connection *models.Connection) {
	if request == nil || connection == nil {
		return
	}

	// Notify the original sender that their request was accepted
	requestEvent := ConnectionRequestEvent{
		Request: request,
		Action:  "accepted",
	}
	s.hub.BroadcastToUser(request.SenderID, EventConnectionRequest, requestEvent)

	// Notify both users about the new connection
	connectionEvent := ConnectionEvent{
		Connection: connection,
		Action:     "established",
	}
	s.hub.BroadcastToUser(request.SenderID, EventConnectionAccepted, connectionEvent)
	s.hub.BroadcastToUser(request.ReceiverID, EventConnectionAccepted, connectionEvent)
}

// BroadcastConnectionDeclined broadcasts that a connection request was declined
func (s *WebSocketService) BroadcastConnectionDeclined(request *models.ConnectionRequest) {
	if request == nil {
		return
	}

	requestEvent := ConnectionRequestEvent{
		Request: request,
		Action:  "declined",
	}

	s.hub.BroadcastToUser(request.SenderID, EventConnectionRequest, requestEvent)
}

// BroadcastConnectionDropped broadcasts that a connection was dropped
func (s *WebSocketService) BroadcastConnectionDropped(connectionID, userAID, userBID uuid.UUID) {
	connection := &models.Connection{
		ID:      connectionID,
		UserAID: userAID,
		UserBID: userBID,
		Status:  "dropped",
	}

	connectionEvent := ConnectionEvent{
		Connection: connection,
		Action:     "dropped",
	}

	s.hub.BroadcastToUser(userAID, EventConnectionDropped, connectionEvent)
	s.hub.BroadcastToUser(userBID, EventConnectionDropped, connectionEvent)
}

// BroadcastUserStatusChange broadcasts user status changes to their connections
func (s *WebSocketService) BroadcastUserStatusChange(userID uuid.UUID, status string) {

	// This will be handled by the hub's internal logic to broadcast to relevant connections
	s.hub.broadcastUserStatus(userID, status)
}

// GetOnlineUsers returns list of currently online users
func (s *WebSocketService) GetOnlineUsers() []uuid.UUID {
	return s.hub.GetOnlineUsers()
}

// IsUserOnline checks if a specific user is currently online
func (s *WebSocketService) IsUserOnline(userID uuid.UUID) bool {
	return s.hub.IsUserOnline(userID)
}

// GetConnectionUsers returns users currently connected to a specific connection
func (s *WebSocketService) GetConnectionUsers(connectionID uuid.UUID) []uuid.UUID {
	return s.hub.GetConnectionUsers(connectionID)
}

// SendDirectMessage sends a direct message to a specific user (for system notifications)
func (s *WebSocketService) SendDirectMessage(userID uuid.UUID, eventType EventType, data interface{}) {
	s.hub.BroadcastToUser(userID, eventType, data)
}

// SendErrorToUser sends an error message to a specific user
func (s *WebSocketService) SendErrorToUser(userID uuid.UUID, code int, message string) {
	errorEvent := ErrorEvent{
		Code:    code,
		Message: message,
	}
	s.hub.BroadcastToUser(userID, EventError, errorEvent)
}

// Shutdown gracefully shuts down the WebSocket service
func (s *WebSocketService) Shutdown() {
	s.hub.Shutdown()
}
