package connections

import (
	"context"
	"fmt"
	"match-me/ent"
	"match-me/ent/message"
	"strings"
	"time"

	"github.com/google/uuid"
)

type messageRepository struct {
	client *ent.Client
}

func NewMessageRepository(client *ent.Client) MessageRepository {
	return &messageRepository{
		client: client,
	}
}

func (r *messageRepository) CreateTextMessage(ctx context.Context, connectionID, senderID, receiverID uuid.UUID, content string) (*ent.Message, error) {
	msg, err := r.client.Message.Create().
		SetConnectionID(connectionID).
		SetSenderID(senderID).
		SetReceiverID(receiverID).
		SetType(message.TypeText).
		SetContent(content).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create text message: %w", err)
	}
	return msg, nil
}

func (r *messageRepository) CreateMediaMessage(ctx context.Context, connectionID, senderID, receiverID uuid.UUID, mediaURL, mediaType, publicID, txtContent string) (*ent.Message, error) {
	create := r.client.Message.Create().
		SetConnectionID(connectionID).
		SetSenderID(senderID).
		SetReceiverID(receiverID).
		SetType(message.TypeMedia).
		SetMediaURL(mediaURL).
		SetMediaType(mediaType).
		SetMediaPublicID(publicID)

	if txtContent != "" {
		create.SetType(message.TypeMixed)
		create.SetContent(txtContent)
	}

	msg, err := create.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create media message: %w", err)
	}
	return msg, nil
}

func (r *messageRepository) GetMessage(ctx context.Context, messageID uuid.UUID) (*ent.Message, error) {
	msg, err := r.client.Message.Get(ctx, messageID)
	if err != nil {
		return nil, fmt.Errorf("failed to get message: %w", err)
	}
	return msg, nil
}

func (r *messageRepository) UpdateMessage(ctx context.Context, messageID uuid.UUID, content string) (*ent.Message, error) {
	msg, err := r.client.Message.UpdateOneID(messageID).
		SetContent(content).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update message: %w", err)
	}
	return msg, nil
}

func (r *messageRepository) DeleteMessage(ctx context.Context, messageID uuid.UUID) error {
	err := r.client.Message.UpdateOneID(messageID).
		SetIsDeleted(true).
		SetDeletedAt(time.Now()).
		Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete message: %w", err)
	}
	return nil
}

func (r *messageRepository) GetConnectionMessages(ctx context.Context, connectionID uuid.UUID, limit, offset int) ([]*ent.Message, error) {
	query := r.client.Message.Query().
		Where(
			message.And(
				message.ConnectionIDEQ(connectionID),
				message.IsDeletedEQ(false),
			),
		).
		Order(message.ByCreatedAt())

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	messages, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection messages: %w", err)
	}
	return messages, nil
}

func (r *messageRepository) GetConnectionMessagesWithUsers(ctx context.Context, connectionID uuid.UUID, limit, offset int) ([]*ent.Message, error) {
	query := r.client.Message.Query().
		Where(
			message.And(
				message.ConnectionIDEQ(connectionID),
				message.IsDeletedEQ(false),
			),
		).
		WithSender().
		WithReceiver().
		Order(ent.Desc(message.FieldCreatedAt))

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	messages, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection messages with users: %w", err)
	}
	return messages, nil
}

func (r *messageRepository) MarkMessageAsRead(ctx context.Context, messageID uuid.UUID) (*ent.Message, error) {
	msg, err := r.client.Message.UpdateOneID(messageID).
		SetIsRead(true).
		SetReadAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to mark message as read: %w", err)
	}
	return msg, nil
}

func (r *messageRepository) MarkConnectionMessagesAsRead(ctx context.Context, connectionID, userID uuid.UUID) (int, error) {
	count, err := r.client.Message.Update().
		Where(
			message.And(
				message.ConnectionIDEQ(connectionID),
				message.ReceiverIDEQ(userID),
				message.IsReadEQ(false),
				message.IsDeletedEQ(false),
			),
		).
		SetIsRead(true).
		SetReadAt(time.Now()).
		Save(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to mark connection messages as read: %w", err)
	}
	return count, nil
}

func (r *messageRepository) GetUnreadMessagesCount(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := r.client.Message.Query().
		Where(
			message.And(
				message.ReceiverIDEQ(userID),
				message.IsReadEQ(false),
				message.IsDeletedEQ(false),
			),
		).
		Count(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread messages count: %w", err)
	}
	return count, nil
}

func (r *messageRepository) GetUnreadMessagesForConnection(ctx context.Context, connectionID, userID uuid.UUID) ([]*ent.Message, error) {
	messages, err := r.client.Message.Query().
		Where(
			message.And(
				message.ConnectionIDEQ(connectionID),
				message.ReceiverIDEQ(userID),
				message.IsReadEQ(false),
				message.IsDeletedEQ(false),
			),
		).
		WithSender().
		Order(message.ByCreatedAt()).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get unread messages for connection: %w", err)
	}
	return messages, nil
}

func (r *messageRepository) GetUserMessages(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*ent.Message, error) {
	query := r.client.Message.Query().
		Where(
			message.And(
				message.Or(
					message.SenderIDEQ(userID),
					message.ReceiverIDEQ(userID),
				),
				message.IsDeletedEQ(false),
			),
		).
		WithConnection().
		WithSender().
		WithReceiver().
		Order(message.ByCreatedAt())

	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}

	messages, err := query.All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user messages: %w", err)
	}
	return messages, nil
}

func (r *messageRepository) GetMediaMessages(ctx context.Context, connectionID uuid.UUID) ([]*ent.Message, error) {
	messages, err := r.client.Message.Query().
		Where(
			message.And(
				message.ConnectionIDEQ(connectionID),
				message.TypeEQ(message.TypeMedia),
				message.IsDeletedEQ(false),
			),
		).
		WithSender().
		Order(message.ByCreatedAt()).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get media messages: %w", err)
	}
	return messages, nil
}

func (r *messageRepository) SearchMessages(ctx context.Context, connectionID uuid.UUID, query string) ([]*ent.Message, error) {
	searchTerm := strings.ToLower(query)

	messages, err := r.client.Message.Query().
		Where(
			message.And(
				message.ConnectionIDEQ(connectionID),
				message.TypeEQ(message.TypeText),
				message.ContentContainsFold(searchTerm),
				message.IsDeletedEQ(false),
			),
		).
		WithSender().
		Order(message.ByCreatedAt()).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to search messages: %w", err)
	}
	return messages, nil
}
