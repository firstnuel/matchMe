package models

import (
	"time"

	"github.com/google/uuid"
)

// Event represents local events (Couchsurfing-style)
type Event struct {
	ID                  uuid.UUID  `json:"id" db:"id"`
	CreatorID           uuid.UUID  `json:"creator_id" db:"creator_id"`
	Title               string     `json:"title" db:"title"`
	Description         *string    `json:"description,omitempty" db:"description"`
	EventType           *string    `json:"event_type,omitempty" db:"event_type"`
	LocationLat         *float64   `json:"location_lat,omitempty" db:"location_lat"`
	LocationLng         *float64   `json:"location_lng,omitempty" db:"location_lng"`
	LocationName        *string    `json:"location_name,omitempty" db:"location_name"`
	MaxParticipants     *int       `json:"max_participants,omitempty" db:"max_participants"`
	CurrentParticipants int        `json:"current_participants" db:"current_participants"`
	StartTime           *time.Time `json:"start_time,omitempty" db:"start_time"`
	EndTime             *time.Time `json:"end_time,omitempty" db:"end_time"`
	IsActive            bool       `json:"is_active" db:"is_active"`
	CreatedAt           time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time  `json:"updated_at" db:"updated_at"`

	// Computed field
	Distance *float64 `json:"distance,omitempty" db:"-"`
}

// EventParticipant represents event participation
type EventParticipant struct {
	EventID  uuid.UUID `json:"event_id" db:"event_id"`
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	Status   string    `json:"status" db:"status"` // 'interested', 'confirmed', 'declined'
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}
