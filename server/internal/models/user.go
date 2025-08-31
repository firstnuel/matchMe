package models

import (
	"match-me/ent"
	"match-me/ent/schema"

	"github.com/google/uuid"
)

// AccessLevel defines the type of data to return
type AccessLevel int

const (
	AccessLevelBasic   AccessLevel = iota // For /users/{id} - name and profile picture
	AccessLevelProfile                    // For /users/{id}/profile - "about me" info
	AccessLevelBio                        // For /users/{id}/bio - recommendation data
	AccessLevelFull                       // For internal/owner use - all data
)

type User struct {
	ID                 uuid.UUID       `json:"id"`
	Email              string          `json:"email,omitempty"`
	FirstName          string          `json:"first_name"`
	LastName           string          `json:"last_name"`
	CreatedAt          *string         `json:"created_at,omitempty"`
	UpdatedAt          *string         `json:"updated_at,omitempty"`
	Age                int             `json:"age,omitempty"`
	AboutMe            *string         `json:"about_me,omitempty"`
	PreferredAgeMin    *int            `json:"preferred_age_min,omitempty"`
	PreferredDistance  *int            `json:"preferred_distance,omitempty"`
	PreferredAgeMax    *int            `json:"preferred_age_max,omitempty"`
	ProfileCompletion  int             `json:"profile_completion,omitempty"`
	Gender             string          `json:"gender,omitempty"`
	PreferredGender    string          `json:"preferred_gender,omitempty"`
	Coordinates        *schema.Point   `json:"coordinates,omitempty"`
	LookingFor         []string        `json:"looking_for,omitempty"`
	Interests          []string        `json:"interests,omitempty"`
	MusicPreferences   []string        `json:"music_preferences,omitempty"`
	FoodPreferences    []string        `json:"food_preferences,omitempty"`
	CommunicationStyle *string         `json:"communication_style,omitempty"`
	Prompts            []schema.Prompt `json:"prompts,omitempty"`
	Photos             []UserPhoto     `json:"photos,omitempty"`
	ProfilePhoto       *string         `json:"profile_photo,omitempty"`
}

type UserPhoto struct {
	ID       uuid.UUID `json:"id"`
	PhotoURL string    `json:"photo_url"`
	Order    int       `json:"order"`
}

func ToUser(entUser *ent.User, accessLevel AccessLevel) *User {
	if entUser == nil {
		return nil
	}

	user := &User{
		ID:        entUser.ID,
		FirstName: entUser.FirstName,
		LastName:  entUser.LastName,
	}

	switch accessLevel {
	case AccessLevelBasic:
		// Only return basic info: ID, name, and profile picture link
		if entUser.Edges.Photos != nil {
			for _, photo := range entUser.Edges.Photos {
				if photo.Order == 1 {
					user.ProfilePhoto = &photo.PhotoURL
					break
				}
			}
		}

	case AccessLevelProfile:
		// Return "about me" type information
		user.Age = entUser.Age
		user.Gender = string(entUser.Gender)

		if entUser.AboutMe != "" {
			user.AboutMe = &entUser.AboutMe
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

	case AccessLevelBio:
		// Return biographical data for recommendations
		user.Age = entUser.Age
		user.Gender = string(entUser.Gender)
		user.PreferredGender = string(entUser.PreferredGender)

		if entUser.AboutMe != "" {
			user.AboutMe = &entUser.AboutMe
		}

		if entUser.LookingFor != nil {
			user.LookingFor = entUser.LookingFor
		}

		if entUser.PreferredAgeMin != 0 {
			user.PreferredAgeMin = &entUser.PreferredAgeMin
		}

		if entUser.PreferredAgeMax != 0 {
			user.PreferredAgeMax = &entUser.PreferredAgeMax
		}

		if entUser.PreferredDistance != 0 {
			user.PreferredDistance = &entUser.PreferredDistance
		}

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
	case AccessLevelFull:
		// Return all data (your original implementation)
		user.Email = entUser.Email
		createdAtStr := entUser.CreatedAt.Format("2006-01-02T15:04:05Z07:00")
		user.CreatedAt = &createdAtStr
		updatedAtStr := entUser.UpdatedAt.Format("2006-01-02T15:04:05Z07:00")
		user.UpdatedAt = &updatedAtStr
		user.Age = entUser.Age
		user.ProfileCompletion = entUser.ProfileCompletion
		user.Gender = string(entUser.Gender)
		user.PreferredGender = string(entUser.PreferredGender)

		user.PreferredAgeMin = &entUser.PreferredAgeMin
		user.PreferredAgeMax = &entUser.PreferredAgeMax
		user.PreferredDistance = &entUser.PreferredDistance

		if entUser.Coordinates != nil {
			user.Coordinates = entUser.Coordinates
		}

		if entUser.AboutMe != "" {
			user.AboutMe = &entUser.AboutMe
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
	}

	return user
}

// UserInteractionStats represents statistics about user interactions
type UserInteractionStats struct {
	DeclinedRequests   int `json:"declined_requests"`
	SkippedProfiles    int `json:"skipped_profiles"`
	DeletedConnections int `json:"deleted_connections"`
	TotalInteractions  int `json:"total_interactions"`
}
