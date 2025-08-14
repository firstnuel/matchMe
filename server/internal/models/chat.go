package models

import (
	"time"

	"github.com/google/uuid"
)

// Conversation represents a chat conversation between two users
type Conversation struct {
	ID        uuid.UUID `json:"id" db:"id"`
	User1ID   uuid.UUID `json:"user1_id" db:"user1_id"`
	User2ID   uuid.UUID `json:"user2_id" db:"user2_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

// Message represents a chat message
type Message struct {
	ID             uuid.UUID `json:"id" db:"id"`
	ConversationID uuid.UUID `json:"conversation_id" db:"conversation_id"`
	SenderID       uuid.UUID `json:"sender_id" db:"sender_id"`
	Content        string    `json:"content" db:"content"`
	MessageType    string    `json:"message_type" db:"message_type"` // 'text', 'image', 'emoji'
	IsRead         bool      `json:"is_read" db:"is_read"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}
