package websocket

import (
	"encoding/json"
	"match-me/internal/models"
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of WebSocket event
type EventType string

const (
	// Message events
	EventMessageNew    EventType = "message_new"
	EventMessageRead   EventType = "message_read"
	EventMessageTyping EventType = "message_typing"

	// User status events
	EventUserOnline       EventType = "user_online"
	EventUserOffline      EventType = "user_offline"
	EventUserAway         EventType = "user_away"
	EventUserStatusChange EventType = "user_status_change"

	// Connection events
	EventConnectionRequest  EventType = "connection_request"
	EventConnectionAccepted EventType = "connection_accepted"
	EventConnectionDropped  EventType = "connection_dropped"

	// System events
	EventError EventType = "error"
	EventPing  EventType = "ping"
	EventPong  EventType = "pong"
)

// WebSocketMessage represents the base structure for WebSocket messages
type WebSocketMessage struct {
	Type      EventType   `json:"type"`
	Data      interface{} `json:"data,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
	MessageID string      `json:"message_id"`
}

// MessageEvent represents a new message event
type MessageEvent struct {
	Message      *models.Message `json:"message"`
	ConnectionID uuid.UUID       `json:"connection_id"`
	SenderID     uuid.UUID       `json:"sender_id"`
	ReceiverID   uuid.UUID       `json:"receiver_id"`
}

// MessageReadEvent represents a message read event
type MessageReadEvent struct {
	MessageID    uuid.UUID `json:"message_id"`
	ConnectionID uuid.UUID `json:"connection_id"`
	ReadBy       uuid.UUID `json:"read_by"`
	ReadAt       time.Time `json:"read_at"`
}

// TypingEvent represents typing indicator event
type TypingEvent struct {
	ConnectionID uuid.UUID `json:"connection_id"`
	UserID       uuid.UUID `json:"user_id"`
	IsTyping     bool      `json:"is_typing"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// UserStatusEvent represents user online/offline status
type UserStatusEvent struct {
	UserID       uuid.UUID `json:"user_id"`
	Status       string    `json:"status"` // "online", "offline", "away"
	LastActivity time.Time `json:"last_activity"`
}

// ConnectionRequestEvent represents connection request events
type ConnectionRequestEvent struct {
	Request *models.ConnectionRequest `json:"request"`
	Action  string                    `json:"action"` // "new", "accepted", "declined"
}

// ConnectionEvent represents connection events
type ConnectionEvent struct {
	Connection *models.Connection `json:"connection"`
	Action     string             `json:"action"` // "established", "dropped"
}

// ErrorEvent represents error events
type ErrorEvent struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// NewWebSocketMessage creates a new WebSocket message with timestamp and ID
func NewWebSocketMessage(eventType EventType, data interface{}) *WebSocketMessage {
	return &WebSocketMessage{
		Type:      eventType,
		Data:      data,
		Timestamp: time.Now(),
		MessageID: uuid.New().String(),
	}
}

// ToJSON converts WebSocketMessage to JSON bytes
func (wsm *WebSocketMessage) ToJSON() ([]byte, error) {
	return json.Marshal(wsm)
}

// FromJSON creates WebSocketMessage from JSON bytes
func FromJSON(data []byte) (*WebSocketMessage, error) {
	var wsm WebSocketMessage
	err := json.Unmarshal(data, &wsm)
	if err != nil {
		return nil, err
	}
	return &wsm, nil
}
