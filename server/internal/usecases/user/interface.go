package user

import (
	"context"
	"match-me/internal/models"
	"match-me/internal/requests"

	"github.com/google/uuid"
)

type UserUsecase interface {
	Register(ctx context.Context, req requests.RegisterUser) (*models.User, string, error)
	Login(ctx context.Context, email, password string) (*models.User, string, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, newPassword string) error
	UpdateUser(ctx context.Context, id uuid.UUID, req *requests.UpdateUser) (*models.User, error)

	GetUserByID(ctx context.Context, id uuid.UUID) (*models.User, error)
	DeleteUser(ctx context.Context, userID uuid.UUID) error
	UploadUserPhotos(ctx context.Context, userID uuid.UUID, files []interface{}) ([]*models.UserPhoto, error)
}
