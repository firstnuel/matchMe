package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Connection struct {
	ent.Schema
}

func (Connection) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),

		field.Enum("status").
			Values("connected", "dropped").
			Comment("Status of the connection between two users"),

		field.UUID("user_a_id", uuid.UUID{}).
			Comment("ID of the first user in the connection"),

		field.UUID("user_b_id", uuid.UUID{}).
			Comment("ID of the second user in the connection"),

		field.Time("connected_at").
			Default(time.Now).
			Comment("Timestamp when the connection was established"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Timestamp when the connection was last updated"),

		field.Time("dropped_at").
			Optional().
			Comment("Timestamp when the connection was dropped (if applicable)"),
	}
}

// Edges of the Connection.
func (Connection) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user_a", User.Type).
			Unique().
			Required().
			Field("user_a_id").
			Comment("Reference to the first user in the connection"),

		edge.To("user_b", User.Type).
			Unique().
			Required().
			Field("user_b_id").
			Comment("Reference to the second user in the connection"),
	}
}

// Indexes of the Connection.
func (Connection) Indexes() []ent.Index {
	return []ent.Index{
		// Index for finding connections by user A
		index.Fields("user_a_id"),

		// Index for finding connections by user B
		index.Fields("user_b_id"),

		// Index for finding connections by status
		index.Fields("status"),

		// Composite index for finding connections between specific users
		index.Fields("user_a_id", "user_b_id").
			Unique(),

		// Index for finding connections by status and user
		index.Fields("user_a_id", "status"),
		index.Fields("user_b_id", "status"),

		// Index for time-based queries
		index.Fields("connected_at"),
		index.Fields("updated_at"),
	}
}
