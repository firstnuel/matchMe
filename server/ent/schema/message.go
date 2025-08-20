package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
)

type Message struct {
	ent.Schema
}

func (Message) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).
			Unique().
			Immutable(),

		field.UUID("connection_id", uuid.UUID{}).
			Comment("ID of the connection this message belongs to"),

		field.UUID("sender_id", uuid.UUID{}).
			Comment("ID of the user who sent the message"),

		field.UUID("receiver_id", uuid.UUID{}).
			Comment("ID of the user who received the message"),

		field.Enum("type").
			Values("text", "media").
			Default("text").
			Comment("Type of message content"),

		field.Text("content").
			Optional().
			Comment("Text content of the message (for text messages)"),

		field.String("media_url").
			Optional().
			Comment("URL of the media file (for media messages)"),

		field.String("media_type").
			Optional().
			Comment("MIME type of the media file (e.g., image/jpeg, video/mp4)"),

		field.String("media_public_id").
			Optional().
			Comment("Public ID for media storage service (e.g., Cloudinary)"),

		field.Bool("is_read").
			Default(false).
			Comment("Whether the message has been read by the receiver"),

		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Comment("Timestamp when the message was sent"),

		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now).
			Comment("Timestamp when the message was last updated"),

		field.Time("read_at").
			Optional().
			Comment("Timestamp when the message was read by the receiver"),

		field.Bool("is_deleted").
			Default(false).
			Comment("Soft delete flag for the message"),

		field.Time("deleted_at").
			Optional().
			Comment("Timestamp when the message was deleted"),
	}
}

// Edges of the Message.
func (Message) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("connection", Connection.Type).
			Unique().
			Required().
			Field("connection_id").
			Comment("Reference to the connection this message belongs to"),

		edge.To("sender", User.Type).
			Unique().
			Required().
			Field("sender_id").
			Comment("Reference to the user who sent the message"),

		edge.To("receiver", User.Type).
			Unique().
			Required().
			Field("receiver_id").
			Comment("Reference to the user who received the message"),
	}
}

// Indexes of the Message.
func (Message) Indexes() []ent.Index {
	return []ent.Index{
		// Index for finding messages by connection (most common query)
		index.Fields("connection_id"),
		
		// Index for finding messages by sender
		index.Fields("sender_id"),
		
		// Index for finding messages by receiver
		index.Fields("receiver_id"),
		
		// Index for finding unread messages
		index.Fields("receiver_id", "is_read"),
		
		// Index for finding messages by connection and timestamp (for pagination)
		index.Fields("connection_id", "created_at"),
		
		// Index for finding messages by type
		index.Fields("type"),
		
		// Index for soft delete queries
		index.Fields("is_deleted"),
		
		// Composite index for connection messages with read status
		index.Fields("connection_id", "is_read"),
		
		// Index for cleanup queries (deleted messages)
		index.Fields("is_deleted", "deleted_at"),
		
		// Index for time-based queries
		index.Fields("created_at"),
		index.Fields("updated_at"),
		
		// Index for media messages
		index.Fields("type", "media_type"),
	}
}