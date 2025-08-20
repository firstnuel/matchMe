package connections

import (
	"context"
	"fmt"
	"match-me/ent"
	"match-me/internal/models"
	"match-me/internal/pkg/cloudinary"
	"match-me/internal/repositories/connections"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type messageUsecase struct {
	messageRepo    connections.MessageRepository
	connectionRepo connections.ConnectionRepository
	cld            cloudinary.Cloudinary
}

func NewMessageUsecase(
	messageRepo connections.MessageRepository,
	connectionRepo connections.ConnectionRepository,
	cld cloudinary.Cloudinary,
) MessageUsecase {
	return &messageUsecase{
		messageRepo:    messageRepo,
		connectionRepo: connectionRepo,
		cld:            cld,
	}
}

func (u *messageUsecase) SendTextMessage(ctx context.Context, senderID uuid.UUID, connectionID uuid.UUID, content string) (*models.Message, error) {
	// Verify the connection exists and user is part of it
	receiverID, err := u.validateConnectionAccess(ctx, senderID, connectionID)
	if err != nil {
		return nil, err
	}

	// Validate content
	if content == "" {
		return nil, fmt.Errorf("message content cannot be empty")
	}

	// Create the message
	entMessage, err := u.messageRepo.CreateTextMessage(ctx, connectionID, senderID, receiverID, content)
	if err != nil {
		return nil, fmt.Errorf("failed to create text message: %w", err)
	}

	return models.ToMessage(entMessage), nil
}

func (u *messageUsecase) SendMediaMessage(ctx context.Context, senderID uuid.UUID, connectionID uuid.UUID, mediaFile interface{}) (*models.Message, error) {
	// Verify the connection exists and user is part of it
	receiverID, err := u.validateConnectionAccess(ctx, senderID, connectionID)
	if err != nil {
		return nil, err
	}

	// Upload image to Cloudinary
	uploadParams := uploader.UploadParams{
		Folder:   "media-photos",
		PublicID: fmt.Sprintf("message_%s_photo", connectionID.String()),
	}

	mediaURL, publicID, err := u.cld.UploadImage(mediaFile, uploadParams)
	if err != nil {
		return nil, fmt.Errorf("failed to upload media to Cloudinary: %w", err)
	}

	// Set mediaType (assuming only images for now)
	mediaType := "image/jpeg"

	// Create the message
	entMessage, err := u.messageRepo.CreateMediaMessage(ctx, connectionID, senderID, receiverID, mediaURL, mediaType, publicID)
	if err != nil {
		return nil, fmt.Errorf("failed to create media message: %w", err)
	}
	return models.ToMessage(entMessage), nil
}

func (u *messageUsecase) GetConnectionMessages(ctx context.Context, userID, connectionID uuid.UUID, limit, offset int) ([]*models.Message, error) {
	// Verify the connection exists and user is part of it
	_, err := u.validateConnectionAccess(ctx, userID, connectionID)
	if err != nil {
		return nil, err
	}

	// Set default limit if not provided
	if limit <= 0 {
		limit = 50
	}

	// Get messages with user details
	entMessages, err := u.messageRepo.GetConnectionMessagesWithUsers(ctx, connectionID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection messages: %w", err)
	}

	return models.ToMessages(entMessages), nil
}

func (u *messageUsecase) MarkMessagesAsRead(ctx context.Context, userID, connectionID uuid.UUID) error {
	// Verify the connection exists and user is part of it
	_, err := u.validateConnectionAccess(ctx, userID, connectionID)
	if err != nil {
		return err
	}

	// Mark all unread messages in this connection as read for this user
	_, err = u.messageRepo.MarkConnectionMessagesAsRead(ctx, connectionID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	return nil
}

func (u *messageUsecase) GetUnreadCount(ctx context.Context, userID uuid.UUID) (int, error) {
	count, err := u.messageRepo.GetUnreadMessagesCount(ctx, userID)
	if err != nil {
		return 0, fmt.Errorf("failed to get unread message count: %w", err)
	}

	return count, nil
}

// validateConnectionAccess verifies that the connection exists, is active, and the user is part of it
// Returns the ID of the other user in the connection
func (u *messageUsecase) validateConnectionAccess(ctx context.Context, userID, connectionID uuid.UUID) (uuid.UUID, error) {
	// Get the connection
	connection, err := u.connectionRepo.GetConnection(ctx, connectionID)
	if err != nil {
		if ent.IsNotFound(err) {
			return uuid.Nil, fmt.Errorf("connection not found")
		}
		return uuid.Nil, fmt.Errorf("failed to get connection: %w", err)
	}

	// Verify connection is active
	if connection.Status != "connected" {
		return uuid.Nil, fmt.Errorf("connection is not active")
	}

	// Verify user is part of this connection and get the other user's ID
	var otherUserID uuid.UUID
	if connection.UserAID == userID {
		otherUserID = connection.UserBID
	} else if connection.UserBID == userID {
		otherUserID = connection.UserAID
	} else {
		return uuid.Nil, fmt.Errorf("unauthorized: user is not part of this connection")
	}

	return otherUserID, nil
}
