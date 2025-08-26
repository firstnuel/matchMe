package interactions

import (
	"context"
	"fmt"
	"match-me/internal/models"
	"match-me/internal/repositories/interactions"
	"time"

	"github.com/google/uuid"
)

// Default expiration periods (hardcoded policies)
const (
	SkippedProfileExpiration    = 30 * 24 * time.Hour  // 30 days
	DeclinedRequestExpiration   = 90 * 24 * time.Hour  // 90 days
	DeletedConnectionExpiration = 365 * 24 * time.Hour // 1 year
)

// Interaction type constants
const (
	InteractionTypeDeclinedRequest   = "declined_request"
	InteractionTypeSkippedProfile    = "skipped_profile"
	InteractionTypeDeletedConnection = "deleted_connection"
)

type userInteractionUsecase struct {
	interactionRepo interactions.UserInteractionRepository
}

func NewUserInteractionUsecase(interactionRepo interactions.UserInteractionRepository) UserInteractionUsecase {
	return &userInteractionUsecase{
		interactionRepo: interactionRepo,
	}
}

func (u *userInteractionUsecase) RecordDeclinedRequest(ctx context.Context, userID, targetUserID uuid.UUID) error {
	expiresAt := time.Now().Add(DeclinedRequestExpiration)
	metadata := map[string]interface{}{
		"reason": "declined_connection_request",
	}

	_, err := u.interactionRepo.RecordInteraction(ctx, userID, targetUserID, InteractionTypeDeclinedRequest, &expiresAt, metadata)
	if err != nil {
		return fmt.Errorf("failed to record declined request: %w", err)
	}

	return nil
}

func (u *userInteractionUsecase) RecordSkippedProfile(ctx context.Context, userID, targetUserID uuid.UUID) error {
	expiresAt := time.Now().Add(SkippedProfileExpiration)
	metadata := map[string]interface{}{
		"reason": "skipped_during_recommendations",
	}

	_, err := u.interactionRepo.RecordInteraction(ctx, userID, targetUserID, InteractionTypeSkippedProfile, &expiresAt, metadata)
	if err != nil {
		return fmt.Errorf("failed to record skipped profile: %w", err)
	}

	return nil
}

func (u *userInteractionUsecase) RecordDeletedConnection(ctx context.Context, userID, targetUserID uuid.UUID) error {
	expiresAt := time.Now().Add(DeletedConnectionExpiration)
	metadata := map[string]interface{}{
		"reason": "deleted_connection",
	}

	_, err := u.interactionRepo.RecordInteraction(ctx, userID, targetUserID, InteractionTypeDeletedConnection, &expiresAt, metadata)
	if err != nil {
		return fmt.Errorf("failed to record deleted connection: %w", err)
	}

	return nil
}

func (u *userInteractionUsecase) GetExcludedUserIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	interactionTypes := []string{
		InteractionTypeDeclinedRequest,
		InteractionTypeSkippedProfile,
		InteractionTypeDeletedConnection,
	}

	excludedIDs, err := u.interactionRepo.GetInteractedUserIDs(ctx, userID, interactionTypes)
	if err != nil {
		return nil, fmt.Errorf("failed to get excluded user IDs: %w", err)
	}

	return excludedIDs, nil
}

func (u *userInteractionUsecase) ShouldExcludeFromRecommendations(ctx context.Context, userID, targetUserID uuid.UUID) (bool, error) {
	// Check for any active interaction that should exclude this user
	interactionTypes := []string{
		InteractionTypeDeclinedRequest,
		InteractionTypeSkippedProfile,
		InteractionTypeDeletedConnection,
	}

	for _, interactionType := range interactionTypes {
		hasInteraction, err := u.interactionRepo.HasActiveInteraction(ctx, userID, targetUserID, interactionType)
		if err != nil {
			return false, fmt.Errorf("failed to check interaction for type %s: %w", interactionType, err)
		}

		if hasInteraction {
			return true, nil
		}
	}

	return false, nil
}

func (u *userInteractionUsecase) ResetSkippedProfiles(ctx context.Context, userID uuid.UUID) error {
	// Get all skipped profile interactions for the user
	skippedInteractions, err := u.interactionRepo.GetActiveInteractionsByType(ctx, userID, InteractionTypeSkippedProfile)
	if err != nil {
		return fmt.Errorf("failed to get skipped interactions: %w", err)
	}

	// Remove each skipped interaction
	for _, interaction := range skippedInteractions {
		err = u.interactionRepo.RemoveInteraction(ctx, userID, interaction.TargetUserID, InteractionTypeSkippedProfile)
		if err != nil {
			return fmt.Errorf("failed to remove skipped interaction: %w", err)
		}
	}

	return nil
}

func (u *userInteractionUsecase) GetInteractionStats(ctx context.Context, userID uuid.UUID) (*models.UserInteractionStats, error) {
	// Get all active interactions for the user
	interactions, err := u.interactionRepo.GetActiveInteractions(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get interactions: %w", err)
	}

	stats := &models.UserInteractionStats{}

	for _, interaction := range interactions {
		switch string(interaction.InteractionType) {
		case InteractionTypeDeclinedRequest:
			stats.DeclinedRequests++
		case InteractionTypeSkippedProfile:
			stats.SkippedProfiles++
		case InteractionTypeDeletedConnection:
			stats.DeletedConnections++
		}
		stats.TotalInteractions++
	}

	return stats, nil
}
