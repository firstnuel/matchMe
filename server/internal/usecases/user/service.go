package user

import (
	"context"
	"errors"
	"fmt"
	"time"

	"match-me/internal/models"
	"match-me/internal/pkg/cloudinary"
	"match-me/internal/pkg/jwt"
	"match-me/internal/repositories/user"
	"match-me/internal/requests"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type userUsecase struct {
	userRepo  user.UserRepository
	jwtSecret string
	cld       cloudinary.Cloudinary
}

func NewUserUsecase(userRepo user.UserRepository, jwtSecret string, cld cloudinary.Cloudinary) UserUsecase {
	return &userUsecase{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
		cld:       cld,
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
		return nil, "", fmt.Errorf("authentication failed: %w", err)
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
