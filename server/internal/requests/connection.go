package requests

import (
	"time"

	"github.com/google/uuid"
)

// Connection represents a connection between two users
type Connection struct {
	ID        uuid.UUID `json:"id" db:"id"`
	User1ID   uuid.UUID `json:"user1_id" db:"user1_id"`
	User2ID   uuid.UUID `json:"user2_id" db:"user2_id"`
	Status    string    `json:"status" db:"status"` // 'pending', 'accepted', 'rejected'
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// ConnectionMode represents different connection types (dating, BFF, networking, events)
type ConnectionMode struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Color       string    `json:"color" db:"color"` // Hex color code
	Icon        *string   `json:"icon,omitempty" db:"icon"`
	Description *string   `json:"description,omitempty" db:"description"`
	IsActive    bool      `json:"is_active" db:"is_active"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// UserInteraction represents user interactions (like/pass)
type UserInteraction struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	TargetUserID    uuid.UUID `json:"target_user_id"`
	InteractionType string    `json:"interaction_type"`
	CreatedAt       time.Time `json:"created_at"`
}
