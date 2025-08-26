package interactions

import (
	"context"
	"fmt"
	"match-me/ent"
	"match-me/ent/userinteraction"
	"time"

	"github.com/google/uuid"
)

type userInteractionRepository struct {
	client *ent.Client
}

func NewUserInteractionRepository(client *ent.Client) UserInteractionRepository {
	return &userInteractionRepository{
		client: client,
	}
}

func (r *userInteractionRepository) RecordInteraction(ctx context.Context, userID, targetUserID uuid.UUID, interactionType string, expiresAt *time.Time, metadata map[string]interface{}) (*ent.UserInteraction, error) {
	mutation := r.client.UserInteraction.Create().
		SetUserID(userID).
		SetTargetUserID(targetUserID).
		SetInteractionType(userinteraction.InteractionType(interactionType))

	if expiresAt != nil {
		mutation = mutation.SetExpiresAt(*expiresAt)
	}

	if metadata != nil {
		mutation = mutation.SetMetadata(metadata)
	}

	// Save the interaction
	interaction, err := mutation.Save(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to record interaction: %w", err)
	}

	return interaction, nil
}

func (r *userInteractionRepository) GetActiveInteractions(ctx context.Context, userID uuid.UUID) ([]*ent.UserInteraction, error) {
	interactions, err := r.client.UserInteraction.Query().
		Where(
			userinteraction.UserID(userID),
			userinteraction.Or(
				userinteraction.ExpiresAtIsNil(),
				userinteraction.ExpiresAtGT(time.Now()),
			),
		).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get active interactions: %w", err)
	}

	return interactions, nil
}

func (r *userInteractionRepository) GetActiveInteractionsByType(ctx context.Context, userID uuid.UUID, interactionType string) ([]*ent.UserInteraction, error) {
	interactions, err := r.client.UserInteraction.Query().
		Where(
			userinteraction.UserID(userID),
			userinteraction.InteractionTypeEQ(userinteraction.InteractionType(interactionType)),
			userinteraction.Or(
				userinteraction.ExpiresAtIsNil(),
				userinteraction.ExpiresAtGT(time.Now()),
			),
		).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get active interactions by type: %w", err)
	}

	return interactions, nil
}

func (r *userInteractionRepository) GetInteractedUserIDs(ctx context.Context, userID uuid.UUID, interactionTypes []string) ([]uuid.UUID, error) {
	if len(interactionTypes) == 0 {
		return []uuid.UUID{}, nil
	}

	// Convert string slice to enum slice
	enumTypes := make([]userinteraction.InteractionType, len(interactionTypes))
	for i, t := range interactionTypes {
		enumTypes[i] = userinteraction.InteractionType(t)
	}

	interactions, err := r.client.UserInteraction.Query().
		Where(
			userinteraction.UserID(userID),
			userinteraction.InteractionTypeIn(enumTypes...),
			userinteraction.Or(
				userinteraction.ExpiresAtIsNil(),
				userinteraction.ExpiresAtGT(time.Now()),
			),
		).
		Select(userinteraction.FieldTargetUserID).
		All(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get interacted user IDs: %w", err)
	}

	userIDs := make([]uuid.UUID, len(interactions))
	for i, interaction := range interactions {
		userIDs[i] = interaction.TargetUserID
	}

	return userIDs, nil
}

func (r *userInteractionRepository) HasActiveInteraction(ctx context.Context, userID, targetUserID uuid.UUID, interactionType string) (bool, error) {
	count, err := r.client.UserInteraction.Query().
		Where(
			userinteraction.UserID(userID),
			userinteraction.TargetUserID(targetUserID),
			userinteraction.InteractionTypeEQ(userinteraction.InteractionType(interactionType)),
			userinteraction.Or(
				userinteraction.ExpiresAtIsNil(),
				userinteraction.ExpiresAtGT(time.Now()),
			),
		).
		Count(ctx)

	if err != nil {
		return false, fmt.Errorf("failed to check active interaction: %w", err)
	}

	return count > 0, nil
}

func (r *userInteractionRepository) RemoveInteraction(ctx context.Context, userID, targetUserID uuid.UUID, interactionType string) error {
	_, err := r.client.UserInteraction.Delete().
		Where(
			userinteraction.UserID(userID),
			userinteraction.TargetUserID(targetUserID),
			userinteraction.InteractionTypeEQ(userinteraction.InteractionType(interactionType)),
		).
		Exec(ctx)

	if err != nil {
		return fmt.Errorf("failed to remove interaction: %w", err)
	}

	return nil
}

func (r *userInteractionRepository) CleanupExpiredInteractions(ctx context.Context) (int, error) {
	count, err := r.client.UserInteraction.Delete().
		Where(
			userinteraction.ExpiresAtNotNil(),
			userinteraction.ExpiresAtLTE(time.Now()),
		).
		Exec(ctx)

	if err != nil {
		return 0, fmt.Errorf("failed to cleanup expired interactions: %w", err)
	}

	return count, nil
}

func (r *userInteractionRepository) GetInteraction(ctx context.Context, interactionID uuid.UUID) (*ent.UserInteraction, error) {
	interaction, err := r.client.UserInteraction.Query().
		Where(userinteraction.ID(interactionID)).
		WithUser().
		WithTargetUser().
		Only(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to get interaction: %w", err)
	}

	return interaction, nil
}