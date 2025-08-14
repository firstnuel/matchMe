package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

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
			Sensitive(),
		field.String("first_name").
			NotEmpty().
			MaxLen(50),
		field.String("username").
			NotEmpty().
			MinLen(3).
			MaxLen(30).
			Unique(),
		field.Time("created_at").
			Default(time.Now).
			Immutable(),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Bool("is_online").
			Default(false),
		field.Int("age").
			Min(18).
			Max(100),
		field.String("gender").
			NotEmpty(),
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
		field.JSON("prompts", []map[string]string{}).
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
