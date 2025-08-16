package hooks

import (
	"context"
	"fmt"
	"match-me/ent"
	"match-me/ent/user"

	"github.com/google/uuid"
)

// ProfileCompletionHook calculates and updates the profile completion percentage
// based on filled fields and photo requirements
func ProfileCompletionHook() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			userMutation, ok := m.(*ent.UserMutation)
			if !ok {
				return next.Mutate(ctx, m)
			}

			// Execute the mutation first
			result, err := next.Mutate(ctx, m)
			if err != nil {
				return result, err
			}

			// Get the user ID to calculate completion
			var userID uuid.UUID
			var hasUserID bool
			if id, exists := userMutation.ID(); exists {
				userID = id
				hasUserID = true
			} else if result != nil {
				if u, ok := result.(*ent.User); ok {
					userID = u.ID
					hasUserID = true
				}
			}

			if !hasUserID {
				return result, nil
			}

			// Calculate and update profile completion
			if err := calculateAndUpdateProfileCompletion(ctx, userMutation.Client(), userID); err != nil {
				// Log error but don't fail the mutation
				fmt.Printf("Error calculating profile completion: %v\n", err)
			}

			return result, nil
		})
	}
}

// calculateAndUpdateProfileCompletion calculates the profile completion percentage
func calculateAndUpdateProfileCompletion(ctx context.Context, client *ent.Client, userID uuid.UUID) error {
	// Fetch the user with photos
	u, err := client.User.Query().
		Where(user.ID(userID)).
		WithPhotos().
		Only(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch user: %w", err)
	}

	completion := calculateCompletion(u)

	// Update the profile completion if it changed
	if u.ProfileCompletion != completion {
		_, err = client.User.UpdateOneID(u.ID).
			SetProfileCompletion(completion).
			Save(ctx)
		if err != nil {
			return fmt.Errorf("failed to update profile completion: %w", err)
		}
	}

	return nil
}

// calculateCompletion calculates the profile completion percentage
func calculateCompletion(u *ent.User) int {
	totalFields := 0
	filledFields := 0

	// Required fields (always count these)
	requiredFields := []bool{
		u.Email != "",
		u.FirstName != "",
		u.LastName != "",
		u.Age > 0,
		u.Gender != "",
	}

	for _, filled := range requiredFields {
		totalFields++
		if filled {
			filledFields++
		}
	}

	// Optional profile fields that contribute to completion
	optionalFields := []bool{
		u.PreferredAgeMin != 0,
		u.PreferredAgeMax != 0,
		u.Coordinates != nil,
		len(u.LookingFor) > 0,
		len(u.Interests) > 0,
		len(u.MusicPreferences) > 0,
		len(u.FoodPreferences) > 0,
		u.CommunicationStyle != "",
		len(u.Prompts) > 0,
	}

	for _, filled := range optionalFields {
		totalFields++
		if filled {
			filledFields++
		}
	}

	// Photos requirement - at least one photo is required
	// This counts as 2 points to emphasize its importance
	totalFields += 2
	if len(u.Edges.Photos) > 0 {
		filledFields += 2
	}

	// Calculate percentage
	if totalFields == 0 {
		return 0
	}

	percentage := (filledFields * 100) / totalFields

	// Ensure it's within the valid range (18-100 as per schema)
	if percentage < 18 {
		return 18
	}
	if percentage > 100 {
		return 100
	}

	return percentage
}

// PhotoCompletionHook specifically handles photo-related changes
// This can be attached to UserPhoto mutations to update profile completion
// when photos are added or removed
func PhotoCompletionHook() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			// Execute the mutation first
			result, err := next.Mutate(ctx, m)
			if err != nil {
				return result, err
			}

			// Only apply to UserPhoto mutations
			photoMutation, ok := m.(*ent.UserPhotoMutation)
			if !ok {
				return result, nil
			}

			// Get the user ID to calculate completion
			var userID uuid.UUID
			var hasUserID bool
			if id, exists := photoMutation.UserID(); exists {
				userID = id
				hasUserID = true
			} else if result != nil {
				if p, ok := result.(*ent.UserPhoto); ok {
					userID = p.UserID
					hasUserID = true
				}
			}

			if !hasUserID {
				return result, nil
			}

			// Update the user's profile completion
			if err := calculateAndUpdateProfileCompletion(ctx, photoMutation.Client(), userID); err != nil {
				// Log error but don't fail the mutation
				fmt.Printf("Error updating profile completion after photo change: %v\n", err)
			}

			return result, nil
		})
	}
}
