package seed

import (
	"context"
	"fmt"
	"log"
	"match-me/ent"
	"match-me/ent/schema"
	"match-me/ent/user"
	"math/rand"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// Seeder handles database population
type Seeder struct {
	client *ent.Client
}

// NewSeeder creates a new seeder instance
func NewSeeder(client *ent.Client) *Seeder {
	return &Seeder{
		client: client,
	}
}

// ResetDatabase drops all data from the database
func (s *Seeder) ResetDatabase(ctx context.Context) error {
	log.Println("Resetting database...")

	// Delete all user photos first (due to foreign key constraint)
	if _, err := s.client.UserPhoto.Delete().Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete user photos: %w", err)
	}

	// Delete all users
	if _, err := s.client.User.Delete().Exec(ctx); err != nil {
		return fmt.Errorf("failed to delete users: %w", err)
	}

	log.Println("Database reset completed successfully")
	return nil
}

// PopulateUsers creates n fake users in the database
func (s *Seeder) PopulateUsers(ctx context.Context, n int) error {
	log.Printf("Starting to populate database with %d users...", n)

	// Sample data pools
	firstNames := []string{
		"Alex", "Jordan", "Taylor", "Casey", "Morgan", "Riley", "Avery", "Cameron", "Blake", "Quinn",
		"Sam", "Jamie", "Reese", "Drew", "Sage", "Parker", "Rowan", "Phoenix", "River", "Skylar",
		"Dakota", "Emery", "Finley", "Hayden", "Kendall", "Lane", "Logan", "Micah", "Nova", "Oakley",
	}

	lastNames := []string{
		"Johnson", "Williams", "Brown", "Jones", "Garcia", "Miller", "Davis", "Rodriguez", "Martinez", "Hernandez",
		"Lopez", "Gonzalez", "Wilson", "Anderson", "Thomas", "Taylor", "Moore", "Jackson", "Martin", "Lee",
		"Perez", "Thompson", "White", "Harris", "Sanchez", "Clark", "Ramirez", "Lewis", "Robinson", "Walker",
	}

	genders := []string{"male", "female", "non_binary"}
	preferredGenders := []string{"male", "female", "non_binary", "all"}

	interests := []string{
		"travel", "music", "movies", "books", "cooking", "fitness", "art", "photography",
		"gaming", "sports", "hiking", "dancing", "yoga", "meditation", "technology",
		"fashion", "food", "wine", "coffee", "pets", "nature", "adventure", "reading",
	}

	musicPreferences := []string{
		"pop", "rock", "jazz", "classical", "hip-hop", "electronic", "country", "folk",
		"blues", "reggae", "indie", "alternative", "r&b", "soul", "funk", "punk",
		"metal", "latin", "world", "ambient",
	}

	foodPreferences := []string{
		"vegetarian", "vegan", "italian", "chinese", "japanese", "mexican", "indian",
		"thai", "french", "mediterranean", "american", "korean", "vietnamese",
		"middle-eastern", "african", "fusion", "seafood", "bbq", "desserts", "street-food",
	}

	lookingFor := []string{"friendship", "relationship", "casual", "networking"}
	communicationStyles := []string{"direct", "thoughtful", "humorous", "analytical", "creative", "empathetic", "casual", "formal", "energetic", "calm"}

	promptQuestions := []string{
		"What's your idea of a perfect weekend?",
		"What's something you're passionate about?",
		"What's the best advice you've ever received?",
		"What's your favorite way to unwind after a long day?",
		"What's a skill you'd love to learn?",
		"What's your biggest goal for the next year?",
		"What's something that always makes you laugh?",
		"What's your favorite type of adventure?",
		"What's something you're really good at?",
		"What's the most interesting place you've visited?",
	}

	// Major US cities for coordinates
	cities := []struct {
		name string
		lat  float64
		lng  float64
	}{
		{"Helsinki", 60.1699, 24.9384},
		{"Espoo", 60.2055, 24.6559},
		{"Tampere", 61.4978, 23.7610},
		{"Vantaa", 60.2934, 25.0378},
		{"Oulu", 65.0121, 25.4651},
		{"Turku", 60.4518, 22.2666},
		{"Jyväskylä", 62.2415, 25.7209},
		{"Lahti", 60.9827, 25.6615},
		{"Kuopio", 62.8924, 27.6770},
		{"Pori", 61.4850, 21.7970},
		{"Kouvola", 60.8681, 26.7042},
		{"Joensuu", 62.6000, 29.7667},
		{"Lappeenranta", 61.0583, 28.1887},
		{"Hämeenlinna", 60.9959, 24.4643},
		{"Seinäjoki", 62.7903, 22.8413},
		{"Rovaniemi", 66.5039, 25.7294},
		{"Mikkeli", 61.6886, 27.2723},
		{"Kotka", 60.4661, 26.9451},
		{"Salo", 60.3833, 23.1333},
		{"Kokkola", 63.8376, 23.1320},
	}

	for i := range n {
		// Generate random user data
		firstName := firstNames[rand.Intn(len(firstNames))]
		lastName := lastNames[rand.Intn(len(lastNames))]
		email := fmt.Sprintf("%s.%s.%d@example.com", firstName, lastName, i)
		age := rand.Intn(43) + 18 // Age between 18-60

		// Hash a default password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
		if err != nil {
			return fmt.Errorf("failed to hash password: %w", err)
		}

		// Random location
		city := cities[rand.Intn(len(cities))]
		coordinates := &schema.Point{
			Latitude:  city.lat + (rand.Float64()-0.5)*0.1, // Add small random offset
			Longitude: city.lng + (rand.Float64()-0.5)*0.1,
		}

		// Random preferences
		userGender := genders[rand.Intn(len(genders))]
		preferredGender := preferredGenders[rand.Intn(len(preferredGenders))]
		prefAgeMin := rand.Intn(10) + 18             // 18-27
		prefAgeMax := prefAgeMin + rand.Intn(20) + 5 // 5-25 years above min
		if prefAgeMax > 60 {
			prefAgeMax = 60
		}
		prefDistance := rand.Intn(50) + 10 // 10-60 km

		// Random bio data
		userInterests := getRandomSubset(interests, randRange(3, 7))    // 3–7 interests
		userMusic := getRandomSubset(musicPreferences, randRange(1, 5)) // 1–5 music preferences
		userFood := getRandomSubset(foodPreferences, randRange(1, 5))   // 1–5 food preferences
		userLookingFor := getRandomSubset(lookingFor, randRange(1, 2))  // 1-2 looking for
		commStyle := communicationStyles[rand.Intn(len(communicationStyles))]

		// Generate random prompts
		prompts := make([]schema.Prompt, 3+rand.Intn(3)) // 3-5 prompts
		usedQuestions := make(map[string]bool)
		for j := 0; j < len(prompts); j++ {
			var question string
			for {
				question = promptQuestions[rand.Intn(len(promptQuestions))]
				if !usedQuestions[question] {
					usedQuestions[question] = true
					break
				}
			}
			prompts[j] = schema.Prompt{
				Question: question,
				Answer:   generateRandomAnswer(),
			}
		}

		aboutMe := generateAboutMe(firstName)

		// Create user directly in database
		createdUser, err := s.client.User.Create().
			SetEmail(email).
			SetPasswordHash(string(hashedPassword)).
			SetFirstName(firstName).
			SetLastName(lastName).
			SetAge(age).
			SetGender(user.Gender(userGender)).
			SetPreferredGender(user.PreferredGender(preferredGender)).
			SetPreferredAgeMin(prefAgeMin).
			SetPreferredAgeMax(prefAgeMax).
			SetPreferredDistance(prefDistance).
			SetCoordinates(coordinates).
			SetAboutMe(aboutMe).
			SetLookingFor(userLookingFor).
			SetInterests(userInterests).
			SetMusicPreferences(userMusic).
			SetFoodPreferences(userFood).
			SetCommunicationStyle(commStyle).
			SetPrompts(prompts).
			SetProfileCompletion(100).
			Save(ctx)

		if err != nil {
			return fmt.Errorf("failed to create user %d: %w", i+1, err)
		}

		// Create placeholder photos for the user
		photoCount := 2 + rand.Intn(4) // 2-5 photos per user
		for j := 0; j < photoCount; j++ {
			photoOrder := j + 1
			placeholderURL := fmt.Sprintf("https://via.placeholder.com/400x600/4A90E2/FFFFFF?text=User+%d+Photo+%d", i+1, photoOrder)
			placeholderPublicID := fmt.Sprintf("seed_user_%d_photo_%d_%d", i+1, photoOrder, time.Now().Unix())

			_, err = s.client.UserPhoto.Create().
				SetPhotoURL(placeholderURL).
				SetPublicID(placeholderPublicID).
				SetOrder(photoOrder).
				SetUserID(createdUser.ID).
				Save(ctx)

			if err != nil {
				return fmt.Errorf("failed to create photo %d for user %d: %w", photoOrder, i+1, err)
			}
		}

		if (i+1)%10 == 0 {
			log.Printf("Created %d/%d users", i+1, n)
		}
	}

	log.Printf("Successfully populated database with %d users", n)
	return nil
}

// getRandomSubset returns a random subset of the given slice
func getRandomSubset(items []string, count int) []string {
	if count >= len(items) {
		return items
	}

	shuffled := make([]string, len(items))
	copy(shuffled, items)

	for i := len(shuffled) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}

	return shuffled[:count]
}

// generateRandomAnswer generates a random answer for prompts
func generateRandomAnswer() string {
	answers := []string{
		"I love exploring new places and trying different cuisines. There's something magical about discovering hidden gems.",
		"Reading a good book with a cup of coffee while listening to jazz music is my perfect evening.",
		"I'm passionate about photography and capturing moments that tell a story.",
		"Hiking in nature helps me clear my mind and appreciate the simple things in life.",
		"I enjoy cooking for friends and family - food brings people together in amazing ways.",
		"Learning new languages has always fascinated me. It opens doors to different cultures.",
		"Working out keeps me energized and motivated throughout the day.",
		"I love live music venues - there's nothing like the energy of a great concert.",
		"Traveling to new destinations and meeting locals always teaches me something new.",
		"Creating art, whether it's painting or digital design, is how I express myself.",
		"Playing board games with friends brings out everyone's competitive and fun side.",
		"Meditation and mindfulness have helped me stay grounded in this busy world.",
		"I'm always up for trying new adventures, from rock climbing to salsa dancing.",
		"Volunteering at local charities gives me a sense of purpose and connection to my community.",
		"Watching foreign films with subtitles has become my favorite way to unwind.",
	}
	return answers[rand.Intn(len(answers))]
}

// generateAboutMe generates a random about me section
func generateAboutMe(firstName string) string {
	templates := []string{
		"Hi! I'm %s, and I love exploring new experiences and meeting interesting people. I believe life is meant to be lived to the fullest!",
		"Hey there! %s here. I'm passionate about making genuine connections and sharing great conversations over coffee or while exploring the city.",
		"Hello! I'm %s, a curious soul who enjoys both quiet evenings and exciting adventures. Looking forward to meeting like-minded people!",
		"Hi, I'm %s! I value authenticity, kindness, and good humor. Let's create some memorable moments together!",
		"Hey! %s here, always ready for the next adventure. I love deep conversations, spontaneous trips, and finding joy in everyday moments.",
	}
	template := templates[rand.Intn(len(templates))]
	return fmt.Sprintf(template, firstName)
}

func randRange(min, max int) int {
	return min + rand.Intn(max-min+1)
}
