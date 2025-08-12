package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents the core user entity
type User struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	Email              string    `json:"email" db:"email"`
	PasswordHash       string    `json:"-" db:"password_hash"` // Never expose password hash in JSON
	FirstName          *string   `json:"first_name,omitempty" db:"first_name"`
	LastName           *string   `json:"last_name,omitempty" db:"last_name"`
	Username           *string   `json:"username,omitempty" db:"username"`
	ProfilePictureURL  *string   `json:"profile_picture,omitempty" db:"profile_picture_url"`
	LocationLat        *float64  `json:"location_lat,omitempty" db:"location_lat"`
	LocationLng        *float64  `json:"location_lng,omitempty" db:"location_lng"`
	MaxDistanceKm      int       `json:"max_distance_km" db:"max_distance_km"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
	LastSeen           time.Time `json:"last_seen" db:"last_seen"`
	IsOnline           bool      `json:"is_online" db:"is_online"`
	
	// Computed fields (not stored in database)
	Distance           *float64 `json:"distance,omitempty" db:"-"`
	CompatibilityScore *float64 `json:"compatibility_score,omitempty" db:"-"`
	
	// Related data (loaded separately)
	Bio     *UserBio     `json:"bio,omitempty" db:"-"`
	Profile *UserProfile `json:"profile,omitempty" db:"-"`
}

// UserBio represents the biographical data for matching (5+ data points as required)
type UserBio struct {
	UserID             uuid.UUID `json:"user_id" db:"user_id"`
	Age                *int      `json:"age,omitempty" db:"age"`
	Gender             *string   `json:"gender,omitempty" db:"gender"`
	LookingFor         []string  `json:"looking_for,omitempty" db:"looking_for"`
	RelationshipGoals  *string   `json:"relationship_goals,omitempty" db:"relationship_goals"`
	Bio                *string   `json:"bio,omitempty" db:"bio"`
	
	// Data Point 1: Interests
	Interests []string `json:"interests,omitempty" db:"interests"`
	
	// Data Point 2: Music Preferences
	MusicPreferences []string `json:"music_preferences,omitempty" db:"music_preferences"`
	
	// Data Point 3: Food Preferences
	FoodPreferences []string `json:"food_preferences,omitempty" db:"food_preferences"`
	
	// Data Point 4: Travel Style
	TravelStyle *string `json:"travel_style,omitempty" db:"travel_style"`
	
	// Data Point 5: Communication Style
	CommunicationStyle *string `json:"communication_style,omitempty" db:"communication_style"`
	
	// Additional data points for enhanced matching
	ValuesBeliefs      []string `json:"values_beliefs,omitempty" db:"values_beliefs"`
	EducationLevel     *string  `json:"education_level,omitempty" db:"education_level"`
	CareerField        *string  `json:"career_field,omitempty" db:"career_field"`
	LifestyleChoices   []string `json:"lifestyle_choices,omitempty" db:"lifestyle_choices"`
	
	// Data Point 6: Long Walks Preference
	LongWalks *bool `json:"long_walks,omitempty" db:"long_walks"`
	
	// Data Point 7: Movie Preferences
	Movies []string `json:"movies,omitempty" db:"movies"`
	
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserProfile represents the profile information (separate from bio for API requirements)
type UserProfile struct {
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	Headline  *string   `json:"headline,omitempty" db:"headline"`
	AboutMe   *string   `json:"about_me,omitempty" db:"about_me"`
	Photos    []string  `json:"photos,omitempty" db:"photos"`
	Prompts   []Prompt  `json:"prompts,omitempty" db:"prompts"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Prompt represents Hinge-style profile prompts
type Prompt struct {
	Question string `json:"question" db:"question"`
	Answer   string `json:"answer" db:"answer"`
}

// Connection represents a connection between two users
type Connection struct {
	ID        uuid.UUID `json:"id" db:"id"`
	User1ID   uuid.UUID `json:"user1_id" db:"user1_id"`
	User2ID   uuid.UUID `json:"user2_id" db:"user2_id"`
	Status    string    `json:"status" db:"status"` // 'pending', 'accepted', 'rejected'
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// UserInteraction represents user interactions (like/pass)
type UserInteraction struct {
	ID             uuid.UUID `json:"id" db:"id"`
	UserID         uuid.UUID `json:"user_id" db:"user_id"`
	TargetUserID   uuid.UUID `json:"target_user_id" db:"target_user_id"`
	InteractionType string   `json:"interaction_type" db:"interaction_type"` // 'like', 'pass', 'super_like'
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

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

// UserModePreference represents user preferences for different connection modes
type UserModePreference struct {
	UserID   uuid.UUID `json:"user_id" db:"user_id"`
	ModeID   uuid.UUID `json:"mode_id" db:"mode_id"`
	IsActive bool      `json:"is_active" db:"is_active"`
	Priority int       `json:"priority" db:"priority"` // 1=primary, 2=secondary, etc.
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

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
	EventID uuid.UUID `json:"event_id" db:"event_id"`
	UserID  uuid.UUID `json:"user_id" db:"user_id"`
	Status  string    `json:"status" db:"status"` // 'interested', 'confirmed', 'declined'
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

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
