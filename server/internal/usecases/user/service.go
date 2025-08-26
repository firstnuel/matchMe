package user

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"time"

	"match-me/ent"
	"match-me/internal/models"
	"match-me/internal/pkg/cloudinary"
	"match-me/internal/pkg/jwt"
	"match-me/internal/repositories/connections"
	"match-me/internal/repositories/user"
	"match-me/internal/requests"
	"match-me/internal/usecases/interactions"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type userUsecase struct {
	userRepo      user.UserRepository
	jwtSecret     string
	cld           cloudinary.Cloudinary
	connRepo      connections.ConnectionRepository
	connReqRepo   connections.ConnectionRequestRepository
	interactionUC interactions.UserInteractionUsecase
}

func NewUserUsecase(userRepo user.UserRepository,
	connRepo connections.ConnectionRepository,
	connReqRepo connections.ConnectionRequestRepository,
	interactionUC interactions.UserInteractionUsecase,
	jwtSecret string, cld cloudinary.Cloudinary) UserUsecase {
	return &userUsecase{
		userRepo:      userRepo,
		connRepo:      connRepo,
		connReqRepo:   connReqRepo,
		interactionUC: interactionUC,
		jwtSecret:     jwtSecret,
		cld:           cld,
	}
}

func (u *userUsecase) Register(ctx context.Context, req requests.RegisterUser) (*models.User, string, error) {
	// Check if user already exists
	existingUser, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err == nil && existingUser != nil {
		return nil, "", errors.New("user with this email already exists")
	}

	// Create user
	entUser, err := u.userRepo.CreateUser(ctx, req)
	if err != nil {
		return nil, "", fmt.Errorf("failed to create user: %w", err)
	}

	// Generate JWT token
	token, err := jwt.GenerateJWTToken(entUser.ID, u.jwtSecret, jwt.PurposeLogin, 24*time.Hour)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	user := models.ToUser(entUser, models.AccessLevelFull)
	return user, token, nil
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (*models.User, string, error) {
	// Authenticate user
	entUser, err := u.userRepo.Authenticate(ctx, email, password)
	if err != nil {
		return nil, "", fmt.Errorf("authentication failed: incorrect email or password")
	}

	// Generate JWT token
	token, err := jwt.GenerateJWTToken(entUser.ID, u.jwtSecret, jwt.PurposeLogin, 24*time.Hour)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate token: %w", err)
	}

	user := models.ToUser(entUser, models.AccessLevelFull)
	return user, token, nil
}

func (u *userUsecase) UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error {
	// Update password in repository
	return u.userRepo.UpdatePassword(ctx, userID, newPassword)
}

func (u *userUsecase) UpdateUser(ctx context.Context, id uuid.UUID, req *requests.UpdateUser) (*models.User, error) {
	// Check if user exists
	_, err := u.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	// Update user
	entUser, err := u.userRepo.UpdateUser(ctx, id, *req)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	user := models.ToUser(entUser, models.AccessLevelFull)
	return user, nil
}

func (u *userUsecase) GetUserByID(ctx context.Context, userID uuid.UUID, accessLevel models.AccessLevel) (*models.User, error) {
	entUser, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	user := models.ToUser(entUser, accessLevel)
	return user, nil
}

func (u *userUsecase) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	// Check if user exists
	_, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Delete user
	err = u.userRepo.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (u *userUsecase) UploadUserPhotos(ctx context.Context, userID uuid.UUID, files []interface{}) ([]*models.UserPhoto, error) {
	// Check if user exists
	_, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if len(files) == 0 {
		return []*models.UserPhoto{}, nil
	}

	var uploadedPhotos []*models.UserPhoto
	var uploadedPublicIDs []string

	// Upload all files
	for i, file := range files {
		order := i + 1 // Start ordering from 1

		// Upload image to Cloudinary
		uploadParams := uploader.UploadParams{
			Folder:   "user-photos",
			PublicID: fmt.Sprintf("user_%s_photo_%d", userID.String(), order),
		}

		photoURL, publicID, err := u.cld.UploadImage(file, uploadParams)
		if err != nil {
			// Rollback: delete all previously uploaded images
			for _, pid := range uploadedPublicIDs {
				_ = u.cld.DeleteImage(pid)
			}
			return nil, fmt.Errorf("failed to upload image %d: %w", order, err)
		}

		uploadedPublicIDs = append(uploadedPublicIDs, publicID)

		// Create photo request
		photoRequest := requests.UserPhoto{
			PhotoUrl: photoURL,
			PID:      publicID,
			Order:    order,
		}

		// Generate photo ID
		photoID := uuid.New()

		// Add photo to repository
		entPhoto, err := u.userRepo.AddPhoto(ctx, photoID, userID, photoRequest)
		if err != nil {
			// Rollback: delete all uploaded images
			for _, pid := range uploadedPublicIDs {
				_ = u.cld.DeleteImage(pid)
			}
			return nil, fmt.Errorf("failed to save photo %d: %w", order, err)
		}

		// Convert to model
		userPhoto := &models.UserPhoto{
			ID:       entPhoto.ID,
			PhotoURL: entPhoto.PhotoURL,
			Order:    entPhoto.Order,
		}

		uploadedPhotos = append(uploadedPhotos, userPhoto)
	}

	return uploadedPhotos, nil
}

func (u *userUsecase) DeleteUserPhoto(ctx context.Context, userID, photoID uuid.UUID) error {
	// Check if user exists
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	// Find the photo to get its public ID for Cloudinary deletion
	var photoToDelete *ent.UserPhoto
	for _, photo := range user.Edges.Photos {
		if photo.ID == photoID {
			photoToDelete = photo
			break
		}
	}

	if photoToDelete == nil {
		return fmt.Errorf("photo not found for user")
	}

	// Delete photo from repository first
	err = u.userRepo.DeletePhoto(ctx, photoID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete photo from database: %w", err)
	}

	go func() {
		// Delete image from Cloudinary if it has a public ID
		// Generate expected public ID based on upload pattern
		err = u.cld.DeleteImage(photoToDelete.PublicID)
		if err != nil {
			// Log the error but don't fail the operation since DB deletion succeeded
			fmt.Printf("Warning: Failed to delete image from Cloudinary (public_id: %s): %v\n", photoToDelete.PublicID, err)
		}
	}()

	return nil
}

func (u *userUsecase) GetRecommendations(ctx context.Context, userID uuid.UUID) ([]string, error) {
	// Fetch users by user preference
	preferredUsers, currentUser, err := u.userRepo.GetUsersByPreference(ctx, userID)
	if err != nil {
		return []string{}, err
	}

	if len(preferredUsers) < 1 {
		return []string{}, nil
	}

	type userRanking struct {
		userID    string
		userScore int
	}

	recommendedIDS := make([]userRanking, len(preferredUsers))

	for i, v := range preferredUsers {
		score := 0
		// Check matches in bio max if 5 per field
		score += countSimilar(currentUser.LookingFor, v.LookingFor)
		score += countSimilar(currentUser.Interests, v.Interests)
		score += countSimilar(currentUser.MusicPreferences, v.MusicPreferences)
		score += countSimilar(currentUser.FoodPreferences, v.FoodPreferences)
		if currentUser.CommunicationStyle == v.CommunicationStyle {
			score += 5
		}
		recommendedIDS[i] = userRanking{
			userID:    v.ID.String(),
			userScore: score,
		}
	}

	// Sort by userScore descending
	sort.Slice(recommendedIDS, func(i, j int) bool {
		return recommendedIDS[i].userScore > recommendedIDS[j].userScore
	})

	// Fetch existing connections
	connections, err := u.connRepo.GetUserConnections(ctx, userID)
	if err != nil {
		log.Printf("failed to get user connections: %v", err) // only log err
	}
	// Create a set of connected user IDs
	connectedUserIDs := make(map[string]struct{})
	for _, conn := range connections {
		if conn.UserAID == userID {
			connectedUserIDs[conn.UserBID.String()] = struct{}{}
		} else if conn.UserBID == userID {
			connectedUserIDs[conn.UserAID.String()] = struct{}{}
		}
	}

	// Fetch excluded user IDs based on interactions
	var excludedUserIDs map[string]struct{}
	if u.interactionUC != nil {
		excludedIDs, err := u.interactionUC.GetExcludedUserIDs(ctx, userID)
		if err != nil {
			log.Printf("failed to get excluded user IDs: %v", err) // only log err
			excludedUserIDs = make(map[string]struct{})
		} else {
			excludedUserIDs = make(map[string]struct{})
			for _, id := range excludedIDs {
				excludedUserIDs[id.String()] = struct{}{}
			}
		}
	} else {
		excludedUserIDs = make(map[string]struct{})
	}

	// Filter out connected, pending, and excluded users, collect top 10
	result := make([]string, 0, 10)
	for _, ranking := range recommendedIDS {
		if len(result) >= 10 {
			break
		}
		
		// Check if user is connected
		if _, connected := connectedUserIDs[ranking.userID]; connected {
			continue
		}
		
		// Check if there's a pending request in either direction
		targetUserID, err := uuid.Parse(ranking.userID)
		if err != nil {
			continue // Skip if userID is invalid
		}
		
		// Check userID -> targetUserID
		req1, _ := u.connReqRepo.GetConnectionRequestBetweenUsers(ctx, userID, targetUserID)
		// Check targetUserID -> userID  
		req2, _ := u.connReqRepo.GetConnectionRequestBetweenUsers(ctx, targetUserID, userID)
		
		if req1 != nil || req2 != nil {
			continue // Skip if there's a pending request in either direction
		}
		
		// Check if excluded by interactions
		if _, excluded := excludedUserIDs[ranking.userID]; excluded {
			continue
		}
		
		result = append(result, ranking.userID)
	}

	return result, nil
}

func (u *userUsecase) SkipRecommendation(ctx context.Context, userID, targetUserID uuid.UUID) error {
	// Validate that users are not the same
	if userID == targetUserID {
		return fmt.Errorf("cannot skip recommendation for yourself")
	}

	// Check if both users exist
	_, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	_, err = u.userRepo.GetByID(ctx, targetUserID)
	if err != nil {
		return fmt.Errorf("target user not found: %w", err)
	}

	// Record the skipped profile interaction
	if u.interactionUC != nil {
		err = u.interactionUC.RecordSkippedProfile(ctx, userID, targetUserID)
		if err != nil {
			return fmt.Errorf("failed to record skipped profile: %w", err)
		}
	}

	return nil
}

func (u *userUsecase) GetDistanceBetweenUsers(ctx context.Context, userAID, userBID uuid.UUID) (float64, error) {
	// Validate that both users exist
	_, err := u.userRepo.GetByID(ctx, userAID)
	if err != nil {
		return 0, fmt.Errorf("user A not found: %w", err)
	}

	_, err = u.userRepo.GetByID(ctx, userBID)
	if err != nil {
		return 0, fmt.Errorf("user B not found: %w", err)
	}

	// Get distance between users from repository
	distance, err := u.userRepo.GetDistanceBetweenUsers(ctx, userAID, userBID)
	if err != nil {
		return 0, fmt.Errorf("failed to get distance between users: %w", err)
	}

	return distance, nil
}
