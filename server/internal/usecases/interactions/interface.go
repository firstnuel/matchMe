package interactions

import (
	"context"
	"match-me/internal/models"

	"github.com/google/uuid"
)

// UserInteractionUsecase handles business logic for user interactions
type UserInteractionUsecase interface {
	// Record when a user declines a connection request
	RecordDeclinedRequest(ctx context.Context, userID, targetUserID uuid.UUID) error

	// Record when a user skips/passes on a recommended profile
	RecordSkippedProfile(ctx context.Context, userID, targetUserID uuid.UUID) error

	// Record when a user deletes a connection
	RecordDeletedConnection(ctx context.Context, userID, targetUserID uuid.UUID) error

	// Get filtered user IDs that should be excluded from recommendations
	GetExcludedUserIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error)

	// Check if a user should be excluded from recommendations
	ShouldExcludeFromRecommendations(ctx context.Context, userID, targetUserID uuid.UUID) (bool, error)

	// Allow user to reset their interaction history (clear skipped profiles)
	ResetSkippedProfiles(ctx context.Context, userID uuid.UUID) error

	// Get user interaction statistics
	GetInteractionStats(ctx context.Context, userID uuid.UUID) (*models.UserInteractionStats, error)
}