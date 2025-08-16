package models

import (
	"time"

	"match-me/ent"
	"match-me/ent/schema"

	"github.com/google/uuid"
)

type User struct {
	ID                 uuid.UUID       `json:"id"`
	Email              string          `json:"email"`
	FirstName          string          `json:"first_name"`
	LastName           string          `json:"last_name"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	Age                int             `json:"age"`
	PreferredAgeMin    *int            `json:"preferred_age_min,omitempty"`
	PreferredAgeMax    *int            `json:"preferred_age_max,omitempty"`
	ProfileCompletion  int             `json:"profile_completion"`
	Gender             string          `json:"gender"`
	PreferredGender    string          `json:"preferred_gender"`
	Coordinates        *schema.Point   `json:"coordinates,omitempty"`
	LookingFor         []string        `json:"looking_for,omitempty"`
	Interests          []string        `json:"interests,omitempty"`
	MusicPreferences   []string        `json:"music_preferences,omitempty"`
	FoodPreferences    []string        `json:"food_preferences,omitempty"`
	CommunicationStyle *string         `json:"communication_style,omitempty"`
	Prompts            []schema.Prompt `json:"prompts,omitempty"`
	Photos             []UserPhoto     `json:"photos,omitempty"`
}

type UserPhoto struct {
	ID       uuid.UUID `json:"id"`
	PhotoURL string    `json:"photo_url"`
	Order    int       `json:"order"`
}

func ToUser(entUser *ent.User) *User {
	if entUser == nil {
		return nil
	}

	user := &User{
		ID:                entUser.ID,
		Email:             entUser.Email,
		FirstName:         entUser.FirstName,
		LastName:          entUser.LastName,
		CreatedAt:         entUser.CreatedAt,
		UpdatedAt:         entUser.UpdatedAt,
		Age:               entUser.Age,
		ProfileCompletion: entUser.ProfileCompletion,
		Gender:            string(entUser.Gender),
		PreferredGender:   string(entUser.PreferredGender),
	}

	user.PreferredAgeMin = &entUser.PreferredAgeMin

	user.PreferredAgeMax = &entUser.PreferredAgeMax

	if entUser.Coordinates != nil {
		user.Coordinates = entUser.Coordinates
	}

	if entUser.LookingFor != nil {
		user.LookingFor = entUser.LookingFor
	}

	if entUser.Interests != nil {
		user.Interests = entUser.Interests
	}

	if entUser.MusicPreferences != nil {
		user.MusicPreferences = entUser.MusicPreferences
	}

	if entUser.FoodPreferences != nil {
		user.FoodPreferences = entUser.FoodPreferences
	}

	if entUser.CommunicationStyle != "" {
		user.CommunicationStyle = &entUser.CommunicationStyle
	}

	if entUser.Prompts != nil {
		user.Prompts = entUser.Prompts
	}

	if entUser.Edges.Photos != nil {
		photos := make([]UserPhoto, len(entUser.Edges.Photos))
		for i, photo := range entUser.Edges.Photos {
			photos[i] = UserPhoto{
				ID:       photo.ID,
				PhotoURL: photo.PhotoURL,
				Order:    photo.Order,
			}
		}
		user.Photos = photos
	}

	return user
}
