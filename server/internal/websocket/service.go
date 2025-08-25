package websocket

import (
	"log"
	"match-me/internal/models"
	"time" // Make sure time is imported

	"github.com/google/uuid"
)

// WebSocketService provides methods to integrate WebSocket with business logic
type WebSocketService struct {
	chatHub   *ChatHub
	typingHub *TypingHub
	statusHub *StatusHub
}

// NewWebSocketService creates a new WebSocket service
func NewWebSocketService(chatHub *ChatHub, typingHub *TypingHub, statusHub *StatusHub) *WebSocketService {
	return &WebSocketService{
		chatHub:   chatHub,
		typingHub: typingHub,
		statusHub: statusHub,
	}
}

// BroadcastNewMessage broadcasts a new message to connection participants
func (s *WebSocketService) BroadcastNewMessage(message *models.Message) {
	if message == nil {
		log.Printf("❌ BroadcastNewMessage called with nil message")
		return
	}

	// NOTE: Assumes MessageEvent is defined in your package.
	messageEvent := MessageEvent{
		Message:      message,
		ConnectionID: message.ConnectionID,
		SenderID:     message.SenderID,
	}

	// Use the ChatHub to broadcast the message
	s.chatHub.BroadcastMessage(message.ConnectionID, messageEvent, message.SenderID)
}

// BroadcastMessageRead broadcasts message read status to connection participants
func (s *WebSocketService) BroadcastMessageRead(messageID, connectionID, readByUserID uuid.UUID) {
	// NOTE: Assumes MessageReadEvent and EventMessageRead are defined
	readEvent := MessageReadEvent{
		MessageID:    messageID,
		ConnectionID: connectionID,
		ReadBy:       readByUserID,
		ReadAt:       time.Now(),
	}

	// Use ChatHub to broadcast read events within the connection
	s.chatHub.BroadcastEvent(connectionID, EventMessageRead, readEvent, readByUserID)
}

// BroadcastTypingIndicator broadcasts typing status to connection participants
func (s *WebSocketService) BroadcastTypingIndicator(connectionID, userID uuid.UUID, isTyping bool) {
	// NOTE: Assumes TypingEvent is defined.
	typingEvent := TypingEvent{
		ConnectionID: connectionID,
		UserID:       userID,
		IsTyping:     isTyping,
	}

	// Use TypingHub to broadcast typing indicators
	s.typingHub.BroadcastTypingIndicator(connectionID, typingEvent, userID)
}

// BroadcastConnectionRequest broadcasts a new connection request to the receiver
func (s *WebSocketService) BroadcastConnectionRequest(request *models.ConnectionRequest) {
	if request == nil {
		return
	}
	// NOTE: Assumes ConnectionRequestEvent is defined.
	requestEvent := ConnectionRequestEvent{
		Request: request,
		Action:  "new",
	}

	// Use StatusHub to send direct messages to users
	s.statusHub.BroadcastToUser(request.ReceiverID, EventConnectionRequest, requestEvent)
}

// BroadcastConnectionAccepted broadcasts that a connection request was accepted
func (s *WebSocketService) BroadcastConnectionAccepted(request *models.ConnectionRequest, connection *models.Connection) {
	if request == nil || connection == nil {
		return
	}

	requestEvent := ConnectionRequestEvent{
		Request: request,
		Action:  "accepted",
	}
	s.statusHub.BroadcastToUser(request.SenderID, EventConnectionRequest, requestEvent)

	// NOTE: Assumes ConnectionEvent is defined.
	connectionEvent := ConnectionEvent{
		Connection: connection,
		Action:     "established",
	}
	s.statusHub.BroadcastToUser(request.SenderID, EventConnectionAccepted, connectionEvent)
	s.statusHub.BroadcastToUser(request.ReceiverID, EventConnectionAccepted, connectionEvent)
}

// BroadcastConnectionMessagesRead broadcasts that all messages in a connection were read.
func (s *WebSocketService) BroadcastConnectionMessagesRead(connectionID, readByUserID uuid.UUID, messageCount int) {

	readEvent := MessageReadEvent{
		MessageID:    uuid.Nil,
		ConnectionID: connectionID,
		ReadBy:       readByUserID,
		ReadAt:       time.Now(),
	}

	// Use the ChatHub to broadcast the read event to other participants in the connection.
	s.chatHub.BroadcastEvent(connectionID, EventMessageRead, readEvent, readByUserID)
}

// BroadcastConnectionDeclined broadcasts that a connection request was declined.
func (s *WebSocketService) BroadcastConnectionDeclined(request *models.ConnectionRequest) {
	if request == nil {
		return
	}

	// NOTE: Assumes ConnectionRequestEvent and EventConnectionRequest are defined.
	requestEvent := ConnectionRequestEvent{
		Request: request,
		Action:  "declined",
	}

	// Use the StatusHub to send the notification directly to the original sender.
	s.statusHub.BroadcastToUser(request.SenderID, EventConnectionRequest, requestEvent)
}

// BroadcastUserStatusChange broadcasts user status changes to their connections
func (s *WebSocketService) BroadcastUserStatusChange(userID uuid.UUID, status string) {
	s.statusHub.broadcastUserStatus(userID, status)
}

// GetOnlineUsers returns list of currently online users
func (s *WebSocketService) GetOnlineUsers() []uuid.UUID {
	return s.statusHub.GetOnlineUsers()
}

// IsUserOnline checks if a specific user is currently online
func (s *WebSocketService) IsUserOnline(userID uuid.UUID) bool {
	return s.statusHub.IsUserOnline(userID)
}

// GetConnectionUsers returns users currently connected to a specific connection
func (s *WebSocketService) GetConnectionUsers(connectionID uuid.UUID) ([]uuid.UUID, bool) {
	// FIX: Call the new method on the ChatHub.
	return s.chatHub.GetConnectionUsers(connectionID)
}

// SendDirectMessage sends a direct message to a specific user (for system notifications)
func (s *WebSocketService) SendDirectMessage(userID uuid.UUID, eventType EventType, data interface{}) {
	s.statusHub.BroadcastToUser(userID, eventType, data)
}

// SendErrorToUser sends an error message to a specific user
func (s *WebSocketService) SendErrorToUser(userID uuid.UUID, code int, message string) {
	// NOTE: Assumes ErrorEvent is defined.
	errorEvent := ErrorEvent{
		Code:    code,
		Message: message,
	}
	s.statusHub.BroadcastToUser(userID, EventError, errorEvent)
}

// Shutdown gracefully shuts down the WebSocket service
func (s *WebSocketService) Shutdown() {
	s.chatHub.Shutdown()
	s.typingHub.Shutdown()
	s.statusHub.Shutdown()
}
