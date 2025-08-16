package requests

import (
	"fmt"
	"slices"
	"strings"
)

// User represents the core user entity
type RegisterUser struct {
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,min=6"`
	FirstName string `json:"first_name" validate:"required,min=2,max=50"`
	LastName  string `json:"last_name" validate:"required,min=3,max=30"`
	Age       int    `json:"age" validate:"required,min=18,max=100"`
	Gender    string `json:"gender" validate:"required,oneof=male female non_binary prefer_not_to_say"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type UpdatePasswordRequest struct {
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UpdateUser struct {
	FirstName         *string   `json:"first_name" validate:"omitempty,min=2,max=50"`
	LastName          *string   `json:"last_name" validate:"omitempty,min=3,max=30"`
	Age               *int      `json:"age" validate:"omitempty,min=18,max=100"`
	Gender            *string   `json:"gender" validate:"omitempty,oneof=male female non_binary prefer_not_to_say"`
	Location          *Location `json:"location" validate:"omitempty"`
	AboutMe           *string   `json:"about_me" validate:"required,min=10,max=500"`
	Bio               *UserBio  `json:"bio" validate:"omitempty"`
	PreferredAgeMin   *int      `json:"preferred_age_min" validate:"omitempty,min=18,max=100"`
	PreferredAgeMax   *int      `json:"preferred_age_max" validate:"omitempty,min=18,max=100"`
	PreferredGender   *string   `json:"preferred_gender" validate:"omitempty,oneof=male female non_binary all"`
	PreferredDistance *int      `json:"preferred_distance" validate:"omitempty,min=0,max=100"`
}

type UserPhoto struct {
	PhotoUrl string `json:"photo_url" validate:"required,url"`
	Order    int    `json:"order" validate:"required,min=1"`
}

type Location struct {
	Latitude  float64 `json:"latitude" validate:"required"`
	Longitude float64 `json:"longitude" validate:"required"`
}

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

// normalizeSlice normalizes a slice of strings by trimming and converting to lowercase
func normalizeSlice(slice []string) []string {
	n := make([]string, len(slice))
	for i, s := range slice {
		n[i] = strings.ToLower(strings.TrimSpace(s))
	}
	return n
}

// ValidateLocation validates the Location struct
func ValidateLocation(loc Location) error {
	var errors []string
	if loc.Latitude < -90 || loc.Latitude > 90 {
		errors = append(errors, fmt.Sprintf("latitude must be between -90 and 90, got %f", loc.Latitude))
	}
	if loc.Longitude < -180 || loc.Longitude > 180 {
		errors = append(errors, fmt.Sprintf("longitude must be between -180 and 180, got %f", loc.Longitude))
	}
	if len(errors) > 0 {
		return fmt.Errorf("location validation failed: %s", strings.Join(errors, "; "))
	}
	return nil
}

// ValidateUserBio validates UserBio fields against predefined options
func ValidateUserBio(bio UserBio) error {
	var errors []string

	// Normalize all strings
	bio.LookingFor = normalizeSlice(bio.LookingFor)
	bio.Interests = normalizeSlice(bio.Interests)
	bio.MusicPreferences = normalizeSlice(bio.MusicPreferences)
	bio.FoodPreferences = normalizeSlice(bio.FoodPreferences)
	bio.CommunicationStyle = strings.ToLower(strings.TrimSpace(bio.CommunicationStyle))

	// Valid option lists
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
		if !slices.Contains(validInterests, interest) {
			errors = append(errors, fmt.Sprintf("invalid interest: %s", interest))
		}
	}

	// Validate music preferences
	for _, music := range bio.MusicPreferences {
		if !slices.Contains(validMusicPreferences, music) {
			errors = append(errors, fmt.Sprintf("invalid music preference: %s", music))
		}
	}

	// Validate food preferences
	for _, food := range bio.FoodPreferences {
		if !slices.Contains(validFoodPreferences, food) {
			errors = append(errors, fmt.Sprintf("invalid food preference: %s", food))
		}
	}

	// Validate communication style
	if !slices.Contains(validCommunicationStyles, bio.CommunicationStyle) {
		errors = append(errors, fmt.Sprintf("invalid communication style: %s", bio.CommunicationStyle))
	}

	// Validate unique prompt questions
	seenQuestions := make(map[string]bool)
	for i, prompt := range bio.Prompts {
		question := strings.ToLower(strings.TrimSpace(prompt.Question))
		if seenQuestions[question] {
			errors = append(errors, fmt.Sprintf("duplicate prompt question at index %d: %s", i, prompt.Question))
		}
		seenQuestions[question] = true
	}

	return nil
}
