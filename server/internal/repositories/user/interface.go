package user

import (
	"context"
	"match-me/ent"
	"match-me/internal/requests"

	"github.com/google/uuid"
)

// UserRepository defines methods for user authentication and management.
// It provides operations for creating, updating, and authenticating users,
type UserRepository interface {
	// Authentication and user creation
	CreateUser(ctx context.Context, userData requests.RegisterUser) (*ent.User, error)
	Authenticate(ctx context.Context, email, password string) (*ent.User, error)
	UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error

	// User retrieval
	GetByID(ctx context.Context, userID uuid.UUID) (*ent.User, error)
	GetUserByEmail(ctx context.Context, email string) (*ent.User, error)

	// User management
	UpdateUser(ctx context.Context, userID uuid.UUID, userData requests.UpdateUser) (*ent.User, error)
	UpdateUserLocation(ctx context.Context, userID uuid.UUID, lat, lng float64) error
	DeleteUser(ctx context.Context, userID uuid.UUID) error

	// Media management
	AddPhoto(ctx context.Context, photoID, userID uuid.UUID, photo requests.UserPhoto) (*ent.UserPhoto, error)
	DeletePhoto(ctx context.Context, photoID, userID uuid.UUID) error

	// Location specific
	GetUsersByPreference(ctx context.Context, reqUserID uuid.UUID) ([]*ent.User, *ent.User, error)
	GetDistanceBetweenUsers(ctx context.Context, userAID, userBID uuid.UUID) (float64, error)
}
