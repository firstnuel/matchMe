package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// UserPhoto holds the schema definition for the UserPhoto entity.
type UserPhoto struct {
	ent.Schema
}

// Fields of the UserPhoto.
func (UserPhoto) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),
		field.String("photo_url").
			NotEmpty(),
		field.Int("order").
			Min(1),
		field.UUID("user_id", uuid.UUID{}),
	}
}

// Edges of the UserPhoto.
func (UserPhoto) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("photos").
			Field("user_id").
			Required().
			Unique(),
	}
}
