package models

import (
	"match-me/ent"

	"github.com/google/uuid"
)

// Connection represents a connection between two users
type Connection struct {
	ID          uuid.UUID `json:"id"`
	UserAID     uuid.UUID `json:"user_a_id"`
	UserBID     uuid.UUID `json:"user_b_id"`
	Status      string    `json:"status"`
	ConnectedAt string    `json:"connected_at"`

	// User details (when loaded with edges)
	UserA *User `json:"user_a,omitempty"`
	UserB *User `json:"user_b,omitempty"`
}

// ConnectionRequest represents a request to connect between two users
type ConnectionRequest struct {
	ID         uuid.UUID `json:"id"`
	SenderID   uuid.UUID `json:"sender_id"`
	ReceiverID uuid.UUID `json:"receiver_id"`
	Status     string    `json:"status"`
	Message    *string   `json:"message,omitempty"`
	CreatedAt  string    `json:"created_at"`

	// User details (when loaded with edges)
	Sender   *User `json:"sender,omitempty"`
	Receiver *User `json:"receiver,omitempty"`
}

// Message represents a message between connected users
type Message struct {
	ID           uuid.UUID `json:"id"`
	ConnectionID uuid.UUID `json:"connection_id"`
	SenderID     uuid.UUID `json:"sender_id"`
	ReceiverID   uuid.UUID `json:"receiver_id"`
	Type         string    `json:"type"`
	Content      *string   `json:"content,omitempty"`
	MediaURL     *string   `json:"media_url,omitempty"`
	MediaType    *string   `json:"media_type,omitempty"`
	IsRead       bool      `json:"is_read"`
	CreatedAt    string    `json:"created_at"`
	ReadAt       *string   `json:"read_at,omitempty"`

	// User and connection details (when loaded with edges)
	Sender     *User       `json:"sender,omitempty"`
	Receiver   *User       `json:"receiver,omitempty"`
	Connection *Connection `json:"connection,omitempty"`
}

// ToConnection converts an ent.Connection to a models.Connection
func ToConnection(entConnection *ent.Connection) *Connection {
	if entConnection == nil {
		return nil
	}

	connection := &Connection{
		ID:          entConnection.ID,
		UserAID:     entConnection.UserAID,
		UserBID:     entConnection.UserBID,
		Status:      string(entConnection.Status),
		ConnectedAt: entConnection.ConnectedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	// Include user details if loaded
	if entConnection.Edges.UserA != nil {
		connection.UserA = ToUser(entConnection.Edges.UserA, AccessLevelBasic)
	}

	if entConnection.Edges.UserB != nil {
		connection.UserB = ToUser(entConnection.Edges.UserB, AccessLevelBasic)
	}

	return connection
}

// ToConnectionRequest converts an ent.ConnectionRequest to a models.ConnectionRequest
func ToConnectionRequest(entRequest *ent.ConnectionRequest) *ConnectionRequest {
	if entRequest == nil {
		return nil
	}

	request := &ConnectionRequest{
		ID:         entRequest.ID,
		SenderID:   entRequest.SenderID,
		ReceiverID: entRequest.ReceiverID,
		Status:     string(entRequest.Status),
		CreatedAt:  entRequest.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if entRequest.Message != "" {
		request.Message = &entRequest.Message
	}

	// Include user details if loaded
	if entRequest.Edges.Sender != nil {
		request.Sender = ToUser(entRequest.Edges.Sender, AccessLevelBasic)
	}

	if entRequest.Edges.Receiver != nil {
		request.Receiver = ToUser(entRequest.Edges.Receiver, AccessLevelBasic)
	}

	return request
}

// ToConnections converts a slice of ent.Connection to models.Connection
func ToConnections(entConnections []*ent.Connection) []*Connection {
	if entConnections == nil {
		return nil
	}

	connections := make([]*Connection, len(entConnections))
	for i, entConnection := range entConnections {
		connections[i] = ToConnection(entConnection)
	}

	return connections
}

// ToConnectionRequests converts a slice of ent.ConnectionRequest to models.ConnectionRequest
func ToConnectionRequests(entRequests []*ent.ConnectionRequest) []*ConnectionRequest {
	if entRequests == nil {
		return nil
	}

	requests := make([]*ConnectionRequest, len(entRequests))
	for i, entRequest := range entRequests {
		requests[i] = ToConnectionRequest(entRequest)
	}

	return requests
}

// ToMessage converts an ent.Message to a models.Message
func ToMessage(entMessage *ent.Message) *Message {
	if entMessage == nil {
		return nil
	}

	message := &Message{
		ID:           entMessage.ID,
		ConnectionID: entMessage.ConnectionID,
		SenderID:     entMessage.SenderID,
		ReceiverID:   entMessage.ReceiverID,
		Type:         string(entMessage.Type),
		IsRead:       entMessage.IsRead,
		CreatedAt:    entMessage.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if entMessage.Content != "" {
		message.Content = &entMessage.Content
	}

	if entMessage.MediaURL != "" {
		message.MediaURL = &entMessage.MediaURL
	}

	if entMessage.MediaType != "" {
		message.MediaType = &entMessage.MediaType
	}

	if !entMessage.ReadAt.IsZero() {
		readAtStr := entMessage.ReadAt.Format("2006-01-02T15:04:05Z07:00")
		message.ReadAt = &readAtStr
	}

	// Include user details if loaded
	if entMessage.Edges.Sender != nil {
		message.Sender = ToUser(entMessage.Edges.Sender, AccessLevelBasic)
	}

	if entMessage.Edges.Receiver != nil {
		message.Receiver = ToUser(entMessage.Edges.Receiver, AccessLevelBasic)
	}

	// Include connection details if loaded
	if entMessage.Edges.Connection != nil {
		message.Connection = ToConnection(entMessage.Edges.Connection)
	}

	return message
}

// ToMessages converts a slice of ent.Message to models.Message
func ToMessages(entMessages []*ent.Message) []*Message {
	if entMessages == nil {
		return nil
	}

	messages := make([]*Message, len(entMessages))
	for i, entMessage := range entMessages {
		messages[i] = ToMessage(entMessage)
	}

	return messages
}

// ChatListItem represents a chat item in the user's chat list
type ChatListItem struct {
	ConnectionID     uuid.UUID `json:"connection_id"`
	OtherUser        *User     `json:"other_user"`
	LastMessage      *Message  `json:"last_message,omitempty"`
	UnreadCount      int       `json:"unread_count"`
	LastActivity     string    `json:"last_activity"`
	ConnectionStatus string    `json:"connection_status"`
}

// ChatList represents a user's complete chat list
type ChatList struct {
	Chats       []*ChatListItem `json:"chats"`
	TotalChats  int             `json:"total_chats"`
	UnreadTotal int             `json:"unread_total"`
}
