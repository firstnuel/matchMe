package connections

import (
	"context"
	"match-me/internal/models"

	"github.com/google/uuid"
)

// ConnectionUsecase handles business logic for user connections
type ConnectionUsecase interface {
	// Get user's connections
	GetUserConnections(ctx context.Context, userID uuid.UUID) ([]*models.Connection, error)

	// Delete a connection
	DeleteConnection(ctx context.Context, userID, connectionID uuid.UUID) error
}

// ConnectionRequestUsecase handles business logic for connection requests
type ConnectionRequestUsecase interface {
	// Send a connection request
	SendRequest(ctx context.Context, senderID, receiverID uuid.UUID, message string) (*models.ConnectionRequest, error)

	// Get pending requests for user
	GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]*models.ConnectionRequest, error)

	// Accept a connection request
	AcceptRequest(ctx context.Context, userID, requestID uuid.UUID) (*models.Connection, error)

	// Decline a connection request
	DeclineRequest(ctx context.Context, userID, requestID uuid.UUID) error
}

// MessageUsecase handles business logic for messaging between connected users
type MessageUsecase interface {
	// Send a text message
	SendTextMessage(ctx context.Context, senderID uuid.UUID, connectionID uuid.UUID, content string) (*models.Message, error)

	// Send a media message
	SendMediaMessage(ctx context.Context, senderID uuid.UUID, connectionID uuid.UUID, mediaFile interface{}) (*models.Message, error)

	// Get messages for a connection
	GetConnectionMessages(ctx context.Context, userID, connectionID uuid.UUID, limit, offset int) ([]*models.Message, error)

	// Mark messages as read
	MarkMessagesAsRead(ctx context.Context, userID, connectionID uuid.UUID) error

	// Get unread message count
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
}
