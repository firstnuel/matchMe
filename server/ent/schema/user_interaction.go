package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type UserInteraction struct {
	ent.Schema
}

func (UserInteraction) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),

		field.UUID("user_id", uuid.UUID{}).
			Comment("ID of the user performing the action"),

		field.UUID("target_user_id", uuid.UUID{}).
			Comment("ID of the user being acted upon"),

		field.Enum("interaction_type").
			Values("declined_request", "skipped_profile", "deleted_connection").
			Comment("Type of interaction performed"),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Timestamp when the interaction was created"),

		field.Time("expires_at").
			Optional().
			Comment("Optional expiration timestamp for the interaction"),

		field.JSON("metadata", map[string]interface{}{}).
			Optional().
			Comment("Optional additional context for the interaction"),
	}
}

// Edges of the UserInteraction.
func (UserInteraction) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Unique().
			Required().
			Field("user_id").
			Comment("Reference to the user who performed the action"),

		edge.To("target_user", User.Type).
			Unique().
			Required().
			Field("target_user_id").
			Comment("Reference to the user who was acted upon"),
	}
}

// Indexes of the UserInteraction.
func (UserInteraction) Indexes() []ent.Index {
	return []ent.Index{
		// Primary index for finding interactions by user
		index.Fields("user_id"),

		// Index for finding interactions by target user
		index.Fields("target_user_id"),

		// Index for finding interactions by type
		index.Fields("interaction_type"),

		// Composite index for user + interaction type queries
		index.Fields("user_id", "interaction_type"),

		// Index for finding interactions with specific users
		index.Fields("user_id", "target_user_id"),

		// Index for time-based queries and expiration cleanup
		index.Fields("created_at"),
		index.Fields("expires_at"),

		// Index for filtering active (non-expired) interactions
		index.Fields("user_id", "interaction_type", "expires_at"),

		// Index for preventing duplicate interactions of the same type
		index.Fields("user_id", "target_user_id", "interaction_type").
			Unique(),
	}
}