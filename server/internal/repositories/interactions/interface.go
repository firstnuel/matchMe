package interactions

import (
	"context"
	"match-me/ent"
	"time"

	"github.com/google/uuid"
)

// UserInteractionRepository defines methods for managing user interactions
type UserInteractionRepository interface {
	// Record a new user interaction
	RecordInteraction(ctx context.Context, userID, targetUserID uuid.UUID, interactionType string, expiresAt *time.Time, metadata map[string]interface{}) (*ent.UserInteraction, error)

	// Get all active (non-expired) interactions for a user
	GetActiveInteractions(ctx context.Context, userID uuid.UUID) ([]*ent.UserInteraction, error)

	// Get active interactions by type
	GetActiveInteractionsByType(ctx context.Context, userID uuid.UUID, interactionType string) ([]*ent.UserInteraction, error)

	// Get user IDs that have active interactions with the specified user
	GetInteractedUserIDs(ctx context.Context, userID uuid.UUID, interactionTypes []string) ([]uuid.UUID, error)

	// Check if there's an active interaction between two users for a specific type
	HasActiveInteraction(ctx context.Context, userID, targetUserID uuid.UUID, interactionType string) (bool, error)

	// Remove/delete an interaction (for cases like "reset" functionality)
	RemoveInteraction(ctx context.Context, userID, targetUserID uuid.UUID, interactionType string) error

	// Cleanup expired interactions
	CleanupExpiredInteractions(ctx context.Context) (int, error)

	// Get interaction by ID
	GetInteraction(ctx context.Context, interactionID uuid.UUID) (*ent.UserInteraction, error)
}