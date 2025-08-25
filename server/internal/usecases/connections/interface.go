package connections

import (
	"context"
	"match-me/internal/models"

	"github.com/google/uuid"
)

// ConnectionUsecase handles business logic for user connections
type ConnectionUsecase interface {
	GetUserConnections(ctx context.Context, userID uuid.UUID) ([]*models.Connection, error)
	DeleteConnection(ctx context.Context, userID, connectionID uuid.UUID) error
}

// ConnectionRequestUsecase handles business logic for connection requests
type ConnectionRequestUsecase interface {
	SendRequest(ctx context.Context, senderID, receiverID uuid.UUID, message string) (*models.ConnectionRequest, error)
	GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]*models.ConnectionRequest, error)
	AcceptRequest(ctx context.Context, userID, requestID uuid.UUID) (*models.Connection, error)
	DeclineRequest(ctx context.Context, userID, requestID uuid.UUID) error
}

// MessageUsecase handles business logic for messaging between connected users
type MessageUsecase interface {
	SendTextMessage(ctx context.Context, senderID uuid.UUID, connectionID uuid.UUID, content string) (*models.Message, error)
	SendMediaMessage(ctx context.Context, senderID uuid.UUID, connectionID uuid.UUID, mediaFile interface{}) (*models.Message, error)
	GetConnectionMessages(ctx context.Context, userID, connectionID uuid.UUID, limit, offset int) ([]*models.Message, error)
	MarkMessagesAsRead(ctx context.Context, userID, connectionID uuid.UUID) error
	GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error)
	GetChatList(ctx context.Context, userID uuid.UUID) (*models.ChatList, error)
}
