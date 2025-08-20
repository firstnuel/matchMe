package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type ConnectionRequest struct {
	ent.Schema
}

func (ConnectionRequest) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),

		field.UUID("sender_id", uuid.UUID{}).
			Comment("ID of the user sending the connection request"),

		field.UUID("receiver_id", uuid.UUID{}).
			Comment("ID of the user receiving the connection request"),

		field.Enum("status").
			Values("pending", "accepted", "declined", "expired").
			Default("pending").
			Comment("Status of the connection request"),

		field.String("message").
			Optional().
			MaxLen(500).
			Comment("Optional message from sender to receiver"),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Timestamp when the request was created"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Timestamp when the request was last updated"),

		field.Time("responded_at").
			Optional().
			Comment("Timestamp when the request was responded to (accepted/declined)"),

		field.Time("expires_at").
			Optional().
			Comment("Timestamp when the request expires (if applicable)"),
	}
}

// Edges of the ConnectionRequest.
func (ConnectionRequest) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("sender", User.Type).
			Unique().
			Required().
			Field("sender_id").
			Comment("Reference to the user who sent the request"),

		edge.To("receiver", User.Type).
			Unique().
			Required().
			Field("receiver_id").
			Comment("Reference to the user who received the request"),
	}
}

// Indexes of the ConnectionRequest.
func (ConnectionRequest) Indexes() []ent.Index {
	return []ent.Index{
		// Index for finding requests by sender
		index.Fields("sender_id"),

		// Index for finding requests by receiver
		index.Fields("receiver_id"),

		// Index for finding requests by status
		index.Fields("status"),

		// Composite index to prevent duplicate requests between same users
		index.Fields("sender_id", "receiver_id").
			Unique(),

		// Index for finding pending requests by receiver
		index.Fields("receiver_id", "status"),

		// Index for finding sent requests by sender
		index.Fields("sender_id", "status"),

		// Index for time-based queries and cleanup
		index.Fields("created_at"),
		index.Fields("updated_at"),
		index.Fields("expires_at"),

		// Index for finding expired requests
		index.Fields("status", "expires_at"),
	}
}
