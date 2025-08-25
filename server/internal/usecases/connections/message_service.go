package connections

import (
	"context"
	"fmt"
	"match-me/ent"
	"match-me/internal/models"
	"match-me/internal/pkg/cloudinary"
	"match-me/internal/repositories/connections"
	"match-me/internal/websocket"

	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"github.com/google/uuid"
)

type messageUsecase struct {
	messageRepo    connections.MessageRepository
	connectionRepo connections.ConnectionRepository
	cld            cloudinary.Cloudinary
	wsService      *websocket.WebSocketService
}

func NewMessageUsecase(
	messageRepo connections.MessageRepository,
	connectionRepo connections.ConnectionRepository,
	cld cloudinary.Cloudinary,
	wsService *websocket.WebSocketService,
) MessageUsecase {
	return &messageUsecase{
		messageRepo:    messageRepo,
		connectionRepo: connectionRepo,
		cld:            cld,
		wsService:      wsService,
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

	message := models.ToMessage(entMessage)

	// Broadcast the new message via WebSocket
	if u.wsService != nil {
		go func(msg *models.Message) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("🚨 PANIC in BroadcastNewMessage goroutine: %v\n", r)
				}
			}()

			fmt.Printf("🔄 Starting BroadcastNewMessage goroutine for message ID: %s\n", msg.ID)
			u.wsService.BroadcastNewMessage(msg)
			fmt.Printf("✅ Completed BroadcastNewMessage goroutine for message ID: %s\n", msg.ID)
		}(message)
	}

	return message, nil
}

func (u *messageUsecase) SendMediaMessage(ctx context.Context, senderID uuid.UUID, connectionID uuid.UUID, mediaFile interface{}, txtContent string) (*models.Message, error) {
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
	entMessage, err := u.messageRepo.CreateMediaMessage(ctx, connectionID, senderID, receiverID, mediaURL, mediaType, publicID, txtContent)
	if err != nil {
		return nil, fmt.Errorf("failed to create media message: %w", err)
	}

	message := models.ToMessage(entMessage)

	// Broadcast the new message via WebSocket
	if u.wsService != nil {
		go func(msg *models.Message) {
			defer func() {
				if r := recover(); r != nil {
					fmt.Printf("🚨 PANIC in BroadcastNewMessage goroutine: %v\n", r)
				}
			}()

			fmt.Printf("🔄 Starting BroadcastNewMessage goroutine for message ID: %s\n", msg.ID)
			u.wsService.BroadcastNewMessage(msg)
			fmt.Printf("✅ Completed BroadcastNewMessage goroutine for message ID: %s\n", msg.ID)
		}(message)
	}

	return message, nil
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
	count, err := u.messageRepo.MarkConnectionMessagesAsRead(ctx, connectionID, userID)
	if err != nil {
		return fmt.Errorf("failed to mark messages as read: %w", err)
	}

	// Broadcast read status via WebSocket
	if u.wsService != nil && count > 0 {
		go func(connID, usrID uuid.UUID, cnt int) {
			u.wsService.BroadcastConnectionMessagesRead(connID, usrID, cnt)
		}(connectionID, userID, count)
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

func (u *messageUsecase) GetChatList(ctx context.Context, userID uuid.UUID) (*models.ChatList, error) {
	// Get user's connections with user details
	entConnections, err := u.connectionRepo.GetConnectionsWithUsers(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user connections: %w", err)
	}

	if len(entConnections) == 0 {
		return &models.ChatList{
			Chats:       []*models.ChatListItem{},
			TotalChats:  0,
			UnreadTotal: 0,
		}, nil
	}

	// Build chat list items
	chatItems := make([]*models.ChatListItem, 0, len(entConnections))
	totalUnread := 0

	for _, entConnection := range entConnections {
		// Determine the other user in the connection
		var otherUser *models.User
		if entConnection.UserAID == userID {
			otherUser = models.ToUser(entConnection.Edges.UserB, models.AccessLevelBasic)
		} else {
			otherUser = models.ToUser(entConnection.Edges.UserA, models.AccessLevelBasic)
		}

		if otherUser == nil {
			// Skip this connection if we can't determine the other user
			continue
		}

		// Get the latest message for this connection
		latestMessages, err := u.messageRepo.GetConnectionMessages(ctx, entConnection.ID, 1, 0)
		var lastMessage *models.Message
		var lastActivity string

		if err == nil && len(latestMessages) > 0 {
			lastMessage = models.ToMessage(latestMessages[0])
			lastActivity = lastMessage.CreatedAt
		} else {
			// If no messages, use connection created time
			lastActivity = entConnection.ConnectedAt.Format("2006-01-02T15:04:05Z07:00")
		}

		// Get unread count for this connection
		unreadMessages, err := u.messageRepo.GetUnreadMessagesForConnection(ctx, entConnection.ID, userID)
		unreadCount := 0
		if err == nil {
			unreadCount = len(unreadMessages)
		}

		totalUnread += unreadCount

		// Create chat list item
		chatItem := &models.ChatListItem{
			ConnectionID:     entConnection.ID,
			OtherUser:        otherUser,
			LastMessage:      lastMessage,
			UnreadCount:      unreadCount,
			LastActivity:     lastActivity,
			ConnectionStatus: string(entConnection.Status),
		}

		chatItems = append(chatItems, chatItem)
	}

	// Sort chat items by last activity (most recent first)
	// Simple bubble sort for now - could be optimized
	for i := 0; i < len(chatItems); i++ {
		for j := i + 1; j < len(chatItems); j++ {
			if chatItems[j].LastActivity > chatItems[i].LastActivity {
				chatItems[i], chatItems[j] = chatItems[j], chatItems[i]
			}
		}
	}

	return &models.ChatList{
		Chats:       chatItems,
		TotalChats:  len(chatItems),
		UnreadTotal: totalUnread,
	}, nil
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
