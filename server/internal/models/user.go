package models

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
)

// User represents the core user entity
type User struct {
	ID           uuid.UUID   `json:"id" validate:"required"`
	Email        string      `json:"email" validate:"required,email"`
	PasswordHash string      `json:"-" validate:"required,min=6"`
	FirstName    string      `json:"first_name" validate:"required,min=2,max=50"`
	Username     string      `json:"username" validate:"required,min=3,max=30,alphanum"`
	CreatedAt    time.Time   `json:"created_at" validate:"required"`
	UpdatedAt    time.Time   `json:"updated_at" validate:"required"`
	IsOnline     bool        `json:"is_online"`
	Age          int         `json:"age" validate:"required,min=18,max=100"`
	Gender       string      `json:"gender" validate:"required"`
	Bio          UserBio     `json:"bio" validate:"required"`
	Photos       []UserPhoto `json:"photos" validate:"required,min=1,max=6,dive"`
}

type UserPhoto struct {
	PhotoUrl string `json:"photo_url" validate:"required,url"`
	Order    int    `json:"order" validate:"required,min=1"`
}

// UserBio represents biographical data for matching
type UserBio struct {
	LookingFor         []string `json:"looking_for" validate:"required,min=1,dive,oneof=friendship relationship casual networking"`
	Interests          []string `json:"interests" validate:"required,min=1,max=10"`
	MusicPreferences   []string `json:"music_preferences" validate:"required,min=1,max=5"`
	FoodPreferences    []string `json:"food_preferences" validate:"required,min=1,max=5"`
	CommunicationStyle string   `json:"communication_style" validate:"required"`
	Prompts            []Prompt `json:"prompts" validate:"required,min=3,max=5,dive"`
}

// Prompt represents profile prompts
type Prompt struct {
	Question string `json:"question" validate:"required,min=10,max=200"`
	Answer   string `json:"answer" validate:"required,min=5,max=500"`
}

// ValidateUserBio validates UserBio fields against predefined options
func ValidateUserBio(bio UserBio) error {
	validInterests := []string{
		"travel", "music", "movies", "books", "cooking", "fitness", "art", "photography",
		"gaming", "sports", "hiking", "dancing", "yoga", "meditation", "technology",
		"fashion", "food", "wine", "coffee", "pets", "nature", "adventure", "reading",
	}

	validMusicPreferences := []string{
		"pop", "rock", "jazz", "classical", "hip-hop", "electronic", "country", "folk",
		"blues", "reggae", "indie", "alternative", "r&b", "soul", "funk", "punk",
		"metal", "latin", "world", "ambient",
	}

	validFoodPreferences := []string{
		"vegetarian", "vegan", "italian", "chinese", "japanese", "mexican", "indian",
		"thai", "french", "mediterranean", "american", "korean", "vietnamese",
		"middle-eastern", "african", "fusion", "seafood", "bbq", "desserts", "street-food",
	}

	validCommunicationStyles := []string{
		"direct", "thoughtful", "humorous", "analytical", "creative", "empathetic",
		"casual", "formal", "energetic", "calm",
	}

	// Validate interests
	for _, interest := range bio.Interests {
		if !slices.Contains(validInterests, strings.ToLower(interest)) {
			return fmt.Errorf("invalid interest: %s", interest)
		}
	}

	// Validate music preferences
	for _, music := range bio.MusicPreferences {
		if !slices.Contains(validMusicPreferences, strings.ToLower(music)) {
			return fmt.Errorf("invalid music preference: %s", music)
		}
	}

	// Validate food preferences
	for _, food := range bio.FoodPreferences {
		if !slices.Contains(validFoodPreferences, strings.ToLower(food)) {
			return fmt.Errorf("invalid food preference: %s", food)
		}
	}

	// Validate communication style
	if !slices.Contains(validCommunicationStyles, strings.ToLower(bio.CommunicationStyle)) {
		return fmt.Errorf("invalid communication style: %s", bio.CommunicationStyle)
	}

	return nil
}
