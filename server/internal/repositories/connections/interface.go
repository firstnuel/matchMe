package connections

import (
	"context"
	"match-me/ent"

	"github.com/google/uuid"
)

// ConnectionRepository defines methods for managing user connections.
type ConnectionRepository interface {
	// Connection management
	CreateConnection(ctx context.Context, userAID, userBID uuid.UUID) (*ent.Connection, error)
	GetConnection(ctx context.Context, connectionID uuid.UUID) (*ent.Connection, error)
	GetConnectionBetweenUsers(ctx context.Context, userAID, userBID uuid.UUID) (*ent.Connection, error)
	UpdateConnectionStatus(ctx context.Context, connectionID uuid.UUID, status string) (*ent.Connection, error)
	DeleteConnection(ctx context.Context, connectionID uuid.UUID) error

	// User connections queries
	GetUserConnections(ctx context.Context, userID uuid.UUID) ([]*ent.Connection, error)
	GetActiveConnections(ctx context.Context, userID uuid.UUID) ([]*ent.Connection, error)
	GetConnectionsWithUsers(ctx context.Context, userID uuid.UUID) ([]*ent.Connection, error)
	GetConnectedUserIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)
}

// ConnectionRequestRepository defines methods for managing connection requests.
type ConnectionRequestRepository interface {
	// Request management
	CreateConnectionRequest(ctx context.Context, senderID, receiverID uuid.UUID, message string) (*ent.ConnectionRequest, error)
	GetConnectionRequest(ctx context.Context, requestID uuid.UUID) (*ent.ConnectionRequest, error)
	GetConnectionRequestBetweenUsers(ctx context.Context, senderID, receiverID uuid.UUID) (*ent.ConnectionRequest, error)
	UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status string) (*ent.ConnectionRequest, error)
	DeleteConnectionRequest(ctx context.Context, requestID uuid.UUID) error

	// User requests queries
	GetPendingRequestsForUser(ctx context.Context, userID uuid.UUID) ([]*ent.ConnectionRequest, error)
	GetSentRequests(ctx context.Context, userID uuid.UUID) ([]*ent.ConnectionRequest, error)
	GetReceivedRequests(ctx context.Context, userID uuid.UUID) ([]*ent.ConnectionRequest, error)
	
	// Request actions
	AcceptRequest(ctx context.Context, requestID uuid.UUID) (*ent.ConnectionRequest, *ent.Connection, error)
	DeclineRequest(ctx context.Context, requestID uuid.UUID) (*ent.ConnectionRequest, error)
	
	// Cleanup
	ExpireOldRequests(ctx context.Context) (int, error)
}

// MessageRepository defines methods for managing messages between connected users.
type MessageRepository interface {
	// Message management
	CreateTextMessage(ctx context.Context, connectionID, senderID, receiverID uuid.UUID, content string) (*ent.Message, error)
	CreateMediaMessage(ctx context.Context, connectionID, senderID, receiverID uuid.UUID, mediaURL, mediaType, publicID string) (*ent.Message, error)
	GetMessage(ctx context.Context, messageID uuid.UUID) (*ent.Message, error)
	UpdateMessage(ctx context.Context, messageID uuid.UUID, content string) (*ent.Message, error)
	DeleteMessage(ctx context.Context, messageID uuid.UUID) error

	// Connection messages
	GetConnectionMessages(ctx context.Context, connectionID uuid.UUID, limit, offset int) ([]*ent.Message, error)
	GetConnectionMessagesWithUsers(ctx context.Context, connectionID uuid.UUID, limit, offset int) ([]*ent.Message, error)
	
	// Read status
	MarkMessageAsRead(ctx context.Context, messageID uuid.UUID) (*ent.Message, error)
	MarkConnectionMessagesAsRead(ctx context.Context, connectionID, userID uuid.UUID) (int, error)
	GetUnreadMessagesCount(ctx context.Context, userID uuid.UUID) (int, error)
	GetUnreadMessagesForConnection(ctx context.Context, connectionID, userID uuid.UUID) ([]*ent.Message, error)

	// Message queries
	GetUserMessages(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ent.Message, error)
	GetMediaMessages(ctx context.Context, connectionID uuid.UUID) ([]*ent.Message, error)
	SearchMessages(ctx context.Context, connectionID uuid.UUID, query string) ([]*ent.Message, error)
}