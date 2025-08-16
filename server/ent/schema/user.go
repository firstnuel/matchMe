package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

// Prompt represents a structured prompt with question and answer
type Prompt struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),

		field.String("email").
			NotEmpty().
			Unique(),

		field.String("password_hash").
			NotEmpty().
			MinLen(6).
			Sensitive(),

		field.String("first_name").
			NotEmpty().
			MinLen(2).
			MaxLen(50),

		field.String("last_name").
			NotEmpty().
			MinLen(3).
			MaxLen(30),

		field.Time("created_at").
			Default(time.Now).
			Immutable(),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),

		field.Int("age").
			Min(18).
			Max(100),

		field.Int("preferred_age_min").
			Optional().
			Min(18).
			Max(100),

		field.Int("preferred_age_max").
			Optional().
			Min(18).
			Max(100),

		field.Int("profile_completion").
			Optional().
			Min(18).
			Max(100),

		field.Enum("gender").
			Values(
				"male",
				"female",
				"non_binary",
				"prefer_not_to_say",
			).
			Comment("Gender options for user profiles, inclusive of non-binary and opt-out preferences"),

		field.Enum("preferred_gender").
			Values(
				"male",
				"female",
				"non_binary",
				"all",
			).
			Default("all").
			Comment("Preferred gender options for user matches"),

		field.Other("coordinates", &Point{}).
			SchemaType(map[string]string{
				dialect.Postgres: "geography(POINT, 4326)",
			}).
			Optional(),

		field.JSON("looking_for", []string{}).
			Optional(),

		field.JSON("interests", []string{}).
			Optional(),

		field.JSON("music_preferences", []string{}).
			Optional(),

		field.JSON("food_preferences", []string{}).
			Optional(),

		field.String("communication_style").
			Optional(),

		field.JSON("prompts", []Prompt{}).
			Optional(),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("photos", UserPhoto.Type).
			Annotations(
				entsql.OnDelete(entsql.Cascade),
			),
	}
}

func (User) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("coordinates").
			Annotations(entsql.IndexType("GIST")),
	}
}
