package user

import (
	"context"
	"fmt"
	"log"
	"match-me/ent"
	"match-me/ent/schema"
	"match-me/ent/user"
	"match-me/ent/userphoto"
	"match-me/internal/requests"
	"math"

	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type userRepository struct {
	client *ent.Client
}

func NewUserRepository(client *ent.Client) UserRepository {
	return &userRepository{
		client: client,
	}
}
func (r *userRepository) CreateUser(ctx context.Context, userData requests.RegisterUser) (*ent.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userData.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password: %w", err)
	}

	userCreate := r.client.User.Create().
		SetEmail(userData.Email).
		SetPasswordHash(string(hashedPassword)).
		SetFirstName(userData.FirstName).
		SetLastName(userData.LastName).
		SetAge(userData.Age).
		SetGender(user.Gender(userData.Gender))
	createdUser, err := userCreate.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return createdUser, nil
}

func (r *userRepository) Authenticate(ctx context.Context, email, password string) (*ent.User, error) {
	user, err := r.client.User.Query().
		Where(user.Email(email)).
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to find user: %w", err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, fmt.Errorf("invalid password")
	}

	return user, nil
}

func (r *userRepository) UpdatePassword(ctx context.Context, userID uuid.UUID, hashedPassword string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(hashedPassword), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("failed to hash password: %w", err)
	}

	err = r.client.User.UpdateOneID(userID).
		SetPasswordHash(string(hash)).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	return nil
}

func (r *userRepository) GetByID(ctx context.Context, userID uuid.UUID) (*ent.User, error) {
	user, err := r.client.User.Query().
		Where(user.ID(userID)).
		WithPhotos().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (*ent.User, error) {
	user, err := r.client.User.Query().
		Where(user.Email(email)).
		WithPhotos().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return user, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, userID uuid.UUID, userData requests.UpdateUser) (*ent.User, error) {
	update := r.client.User.UpdateOneID(userID)

	// Optional: First Name
	if userData.FirstName != nil && *userData.FirstName != "" {
		update = update.SetFirstName(*userData.FirstName)
	}

	// Optional: Last Name
	if userData.LastName != nil && *userData.LastName != "" {
		update = update.SetLastName(*userData.LastName)
	}

	// Optional: Age
	if userData.Age != nil && *userData.Age != 0 {
		update = update.SetAge(*userData.Age)
	}

	// Optional: Gender
	if userData.Gender != nil && *userData.Gender != "" {
		update = update.SetGender(user.Gender(*userData.Gender))
	}

	// Optional: Location
	if userData.Location != nil {
		coordinates := &schema.Point{
			Longitude: userData.Location.Longitude,
			Latitude:  userData.Location.Latitude,
		}
		update = update.SetCoordinates(coordinates)
	}

	// Optional: Bio
	if userData.Bio != nil {
		if len(userData.Bio.LookingFor) > 0 {
			update = update.SetLookingFor(userData.Bio.LookingFor)
		}
		if len(userData.Bio.Interests) > 0 {
			update = update.SetInterests(userData.Bio.Interests)
		}
		if len(userData.Bio.MusicPreferences) > 0 {
			update = update.SetMusicPreferences(userData.Bio.MusicPreferences)
		}
		if len(userData.Bio.FoodPreferences) > 0 {
			update = update.SetFoodPreferences(userData.Bio.FoodPreferences)
		}
		if userData.Bio.CommunicationStyle != "" {
			update = update.SetCommunicationStyle(userData.Bio.CommunicationStyle)
		}
		if len(userData.Bio.Prompts) > 0 {
			prompts := make([]schema.Prompt, len(userData.Bio.Prompts))
			for i, p := range userData.Bio.Prompts {
				prompts[i] = schema.Prompt{
					Question: p.Question,
					Answer:   p.Answer,
				}
			}
			update = update.SetPrompts(prompts)
		}
	}
	if userData.AboutMe != nil {
		update = update.SetAboutMe(*userData.AboutMe)
	}

	// Optional: Preferred Age Min
	if userData.PreferredAgeMin != nil {
		update = update.SetPreferredAgeMin(*userData.PreferredAgeMin)
	}

	// Optional: Preferred Age Max
	if userData.PreferredAgeMax != nil {
		update = update.SetPreferredAgeMax(*userData.PreferredAgeMax)
	}

	// Optional: Preferred Distance
	if userData.PreferredDistance != nil {
		update = update.SetPreferredDistance(*userData.PreferredDistance)
	}

	// Optional: Preferred Gender
	if userData.PreferredGender != nil {
		update = update.SetPreferredGender(user.PreferredGender(*userData.PreferredGender))
	}

	updatedUser, err := update.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return updatedUser, nil
}

func (r *userRepository) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	err := r.client.User.DeleteOneID(userID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	return nil
}

func (r *userRepository) AddPhoto(ctx context.Context, photoID, userID uuid.UUID, photo requests.UserPhoto) (*ent.UserPhoto, error) {
	userPhoto, err := r.client.UserPhoto.Create().
		SetID(photoID).
		SetPhotoURL(photo.PhotoUrl).
		SetPublicID(photo.PID).
		SetOrder(photo.Order).
		SetUserID(userID).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to add photo: %w", err)
	}

	return userPhoto, nil
}

func (r *userRepository) DeletePhoto(ctx context.Context, photoID, userID uuid.UUID) error {
	_, err := r.client.UserPhoto.Delete().
		Where(
			userphoto.ID(photoID),
			userphoto.UserID(userID),
		).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete photo: %w", err)
	}

	return nil
}

func (r *userRepository) UpdateUserLocation(ctx context.Context, userID uuid.UUID, lat, lng float64) error {
	coordinates := &schema.Point{
		Longitude: lng,
		Latitude:  lat,
	}

	err := r.client.User.UpdateOneID(userID).
		SetCoordinates(coordinates).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to update user location: %w", err)
	}

	return nil
}

func (r *userRepository) GetDistanceBetweenUsers(ctx context.Context, userAID, userBID uuid.UUID) (float64, error) {
	// Get user B's coordinates
	userB, err := r.client.User.Get(ctx, userBID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user B: %w", err)
	}

	if userB.Coordinates == nil {
		return 0, fmt.Errorf("user B has no coordinates")
	}

	// Get user A's coordinates to validate
	userA, err := r.client.User.Get(ctx, userAID)
	if err != nil {
		return 0, fmt.Errorf("failed to get user A: %w", err)
	}

	if userA.Coordinates == nil {
		return 0, fmt.Errorf("user A has no coordinates")
	}

	if math.IsNaN(userA.Coordinates.Longitude) || math.IsNaN(userA.Coordinates.Latitude) ||
		math.IsNaN(userB.Coordinates.Longitude) || math.IsNaN(userB.Coordinates.Latitude) {
		return 0, fmt.Errorf("invalid coordinates: userA(long=%v, lat=%v), userB(long=%v, lat=%v)",
			userA.Coordinates.Longitude, userA.Coordinates.Latitude,
			userB.Coordinates.Longitude, userB.Coordinates.Latitude)
	}

	// Calculate distance using the Haversine formula (simple fallback approach)
	lat1 := userA.Coordinates.Latitude
	lon1 := userA.Coordinates.Longitude
	lat2 := userB.Coordinates.Latitude
	lon2 := userB.Coordinates.Longitude

	const R = 6371 // Earth's radius in kilometers

	dLat := (lat2 - lat1) * math.Pi / 180
	dLon := (lon2 - lon1) * math.Pi / 180

	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180

	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c

	// Distance is returned in kilometers  
	return distance, nil
}

func (r *userRepository) GetUsersByPreference(
	ctx context.Context,
	reqUserID uuid.UUID) ([]*ent.User, *ent.User, error) {

	currentUser, err := r.client.User.Get(ctx, reqUserID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get current user: %w", err)
	}

	if currentUser.ProfileCompletion < 100 {
		return nil, nil, fmt.Errorf("current user profile must be complete")
	}

	if currentUser.Coordinates == nil {
		return nil, nil, fmt.Errorf("current user has no coordinates")
	}
	if math.IsNaN(currentUser.Coordinates.Longitude) || math.IsNaN(currentUser.Coordinates.Latitude) {
		return nil, nil, fmt.Errorf("invalid coordinates: longitude=%v, latitude=%v",
			currentUser.Coordinates.Longitude, currentUser.Coordinates.Latitude)
	}
	if currentUser.PreferredDistance <= 0 {
		return nil, nil, fmt.Errorf("invalid preferred distance: %v", currentUser.PreferredDistance)
	}
	distanceInMeters := float64(currentUser.PreferredDistance) * 1000

	query := r.client.User.Query().Where(
		user.IDNEQ(reqUserID),
		user.PreferredAgeMinGTE(currentUser.PreferredAgeMin),
		user.PreferredAgeMaxLTE(currentUser.PreferredAgeMax),
		user.ProfileCompletionGTE(95),
		func(s *sql.Selector) {
			// Use numbered placeholders to align with Ent's query builder
			s.Where(sql.ExprP(
				"ST_DWithin(coordinates, ST_SetSRID(ST_MakePoint($5, $6), 4326)::geography, $7)",
				currentUser.Coordinates.Longitude, currentUser.Coordinates.Latitude, distanceInMeters,
			))
		},
	)

	if currentUser.PreferredGender != user.PreferredGenderAll {
		query = query.Where(user.GenderEQ(user.Gender(currentUser.PreferredGender)))
	}

	users, err := query.All(ctx)
	if err != nil {
		log.Printf("Query failed: %v", err)
		return nil, nil, fmt.Errorf("failed to get users within range: %w", err)
	}

	return users, currentUser, nil
}
