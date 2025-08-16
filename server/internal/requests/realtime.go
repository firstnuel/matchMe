package requests

import (
	"time"

	"github.com/google/uuid"
)

// UserStatus represents live status indicators
type UserStatus struct {
	UserID        uuid.UUID `json:"user_id" db:"user_id"`
	CurrentModeID uuid.UUID `json:"current_mode_id" db:"current_mode_id"`
	IsAvailable   bool      `json:"is_available" db:"is_available"`
	StatusMessage *string   `json:"status_message,omitempty" db:"status_message"`
	LastActivity  time.Time `json:"last_activity" db:"last_activity"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// TypingIndicator represents real-time typing status
type TypingIndicator struct {
	ConversationID uuid.UUID `json:"conversation_id" db:"conversation_id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	IsTyping       bool      `json:"is_typing" db:"is_typing"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}
