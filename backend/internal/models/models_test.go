package models

import (
	"testing"
	"time"

	"github.com/google/uuid"
)

func TestUserModel(t *testing.T) {
	// Test User struct creation
	user := &User{
		ID:            uuid.New(),
		Email:         "test@example.com",
		PasswordHash:  "hashed_password",
		FirstName:     stringPtr("John"),
		LastName:      stringPtr("Doe"),
		Username:      stringPtr("johndoe"),
		MaxDistanceKm: 50,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
		LastSeen:      time.Now(),
		IsOnline:      false,
	}

	// Verify required fields
	if user.Email == "" {
		t.Error("User email should not be empty")
	}

	if user.PasswordHash == "" {
		t.Error("User password hash should not be empty")
	}

	if user.MaxDistanceKm <= 0 {
		t.Error("User max distance should be positive")
	}

	// Verify UUID is generated
	if user.ID == uuid.Nil {
		t.Error("User ID should be a valid UUID")
	}

	t.Logf("User created successfully: %s (%s)", user.Email, user.ID)
}

func TestUserBioModel(t *testing.T) {
	// Test UserBio struct creation
	bio := &UserBio{
		UserID:             uuid.New(),
		Age:                intPtr(25),
		Gender:             stringPtr("other"),
		LookingFor:         []string{"dating", "friendship"},
		RelationshipGoals:  stringPtr("long-term"),
		Bio:                stringPtr("I love meeting new people!"),
		Interests:          []string{"reading", "traveling", "cooking"},
		MusicPreferences:   []string{"rock", "jazz"},
		FoodPreferences:    []string{"italian", "japanese"},
		TravelStyle:        stringPtr("adventure"),
		CommunicationStyle: stringPtr("direct"),
		LongWalks:          boolPtr(true),
		Movies:             []string{"Inception", "The Matrix"},
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}

	// Verify we have at least 5 biographical data points
	dataPoints := 0
	if len(bio.Interests) > 0 {
		dataPoints++
	}
	if len(bio.MusicPreferences) > 0 {
		dataPoints++
	}
	if len(bio.FoodPreferences) > 0 {
		dataPoints++
	}
	if bio.TravelStyle != nil {
		dataPoints++
	}
	if bio.CommunicationStyle != nil {
		dataPoints++
	}
	if bio.LongWalks != nil {
		dataPoints++
	}
	if len(bio.Movies) > 0 {
		dataPoints++
	}

	if dataPoints < 5 {
		t.Errorf("UserBio should have at least 5 data points, got %d", dataPoints)
	}

	t.Logf("UserBio created successfully with %d data points", dataPoints)
}

func TestUserProfileModel(t *testing.T) {
	// Test UserProfile struct creation
	profile := &UserProfile{
		UserID:   uuid.New(),
		Headline: stringPtr("Adventure seeker and coffee lover"),
		AboutMe:  stringPtr("I'm passionate about exploring new places and meeting interesting people."),
		Photos:   []string{"photo1.jpg", "photo2.jpg"},
		Prompts: []Prompt{
			{Question: "What's your ideal weekend?", Answer: "Hiking and coffee shops!"},
			{Question: "Best travel memory?", Answer: "Getting lost in Tokyo!"},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Verify profile has content
	if profile.Headline == nil || *profile.Headline == "" {
		t.Error("Profile headline should not be empty")
	}

	if profile.AboutMe == nil || *profile.AboutMe == "" {
		t.Error("Profile about me should not be empty")
	}

	if len(profile.Photos) == 0 {
		t.Error("Profile should have at least one photo")
	}

	if len(profile.Prompts) == 0 {
		t.Error("Profile should have at least one prompt")
	}

	t.Logf("UserProfile created successfully with %d photos and %d prompts", 
		len(profile.Photos), len(profile.Prompts))
}

func TestConnectionModel(t *testing.T) {
	// Test Connection struct creation
	user1ID := uuid.New()
	user2ID := uuid.New()

	connection := &Connection{
		ID:        uuid.New(),
		User1ID:   user1ID,
		User2ID:   user2ID,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Verify connection has valid users
	if connection.User1ID == uuid.Nil {
		t.Error("Connection User1ID should not be nil")
	}

	if connection.User2ID == uuid.Nil {
		t.Error("Connection User2ID should not be nil")
	}

	if connection.User1ID == connection.User2ID {
		t.Error("Connection should be between different users")
	}

	// Verify status is valid
	validStatuses := map[string]bool{"pending": true, "accepted": true, "rejected": true}
	if !validStatuses[connection.Status] {
		t.Errorf("Connection status should be valid, got: %s", connection.Status)
	}

	t.Logf("Connection created successfully between users %s and %s", 
		connection.User1ID, connection.User2ID)
}

func TestEventModel(t *testing.T) {
	// Test Event struct creation
	event := &Event{
		ID:                  uuid.New(),
		CreatorID:           uuid.New(),
		Title:               "Coffee Meetup",
		Description:         stringPtr("Let's grab coffee and chat!"),
		EventType:           stringPtr("hangout"),
		LocationLat:         float64Ptr(60.1699),
		LocationLng:         float64Ptr(24.9384),
		LocationName:        stringPtr("Starbucks Downtown"),
		MaxParticipants:     intPtr(10),
		CurrentParticipants: 0,
		StartTime:           timePtr(time.Now().Add(24 * time.Hour)),
		EndTime:             timePtr(time.Now().Add(26 * time.Hour)),
		IsActive:            true,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	// Verify event has required fields
	if event.Title == "" {
		t.Error("Event title should not be empty")
	}

	if event.CreatorID == uuid.Nil {
		t.Error("Event creator ID should not be nil")
	}

	if event.MaxParticipants != nil && *event.MaxParticipants <= 0 {
		t.Error("Event max participants should be positive")
	}

	if event.StartTime != nil && event.EndTime != nil {
		if event.StartTime.After(*event.EndTime) {
			t.Error("Event start time should be before end time")
		}
	}

	t.Logf("Event created successfully: %s", event.Title)
}

// Helper functions for creating pointers
func stringPtr(s string) *string {
	return &s
}

func intPtr(i int) *int {
	return &i
}

func boolPtr(b bool) *bool {
	return &b
}

func float64Ptr(f float64) *float64 {
	return &f
}

func timePtr(t time.Time) *time.Time {
	return &t
}
