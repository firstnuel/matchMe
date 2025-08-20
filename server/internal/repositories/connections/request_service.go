package connections

import (
	"context"
	"fmt"
	"match-me/ent"
	"match-me/ent/connection"
	"match-me/ent/connectionrequest"
	"time"

	"github.com/google/uuid"
)

type connectionRequestRepository struct {
	client *ent.Client
}

func NewConnectionRequestRepository(client *ent.Client) ConnectionRequestRepository {
	return &connectionRequestRepository{
		client: client,
	}
}

func (r *connectionRequestRepository) CreateConnectionRequest(ctx context.Context, senderID, receiverID uuid.UUID, message string) (*ent.ConnectionRequest, error) {
	create := r.client.ConnectionRequest.Create().
		SetSenderID(senderID).
		SetReceiverID(receiverID).
		SetStatus(connectionrequest.StatusPending)

	if message != "" {
		create = create.SetMessage(message)
	}

	request, err := create.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection request: %w", err)
	}
	return request, nil
}

func (r *connectionRequestRepository) GetConnectionRequest(ctx context.Context, requestID uuid.UUID) (*ent.ConnectionRequest, error) {
	request, err := r.client.ConnectionRequest.Get(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection request: %w", err)
	}
	return request, nil
}

func (r *connectionRequestRepository) GetConnectionRequestBetweenUsers(ctx context.Context, senderID, receiverID uuid.UUID) (*ent.ConnectionRequest, error) {
	request, err := r.client.ConnectionRequest.Query().
		Where(
			connectionrequest.And(
				connectionrequest.SenderIDEQ(senderID),
				connectionrequest.ReceiverIDEQ(receiverID),
			),
		).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get connection request between users: %w", err)
	}
	return request, nil
}

func (r *connectionRequestRepository) UpdateRequestStatus(ctx context.Context, requestID uuid.UUID, status string) (*ent.ConnectionRequest, error) {
	update := r.client.ConnectionRequest.UpdateOneID(requestID).
		SetStatus(connectionrequest.Status(status)).
		SetRespondedAt(time.Now())

	request, err := update.Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update request status: %w", err)
	}
	return request, nil
}

func (r *connectionRequestRepository) DeleteConnectionRequest(ctx context.Context, requestID uuid.UUID) error {
	err := r.client.ConnectionRequest.DeleteOneID(requestID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete connection request: %w", err)
	}
	return nil
}

func (r *connectionRequestRepository) GetPendingRequestsForUser(ctx context.Context, userID uuid.UUID) ([]*ent.ConnectionRequest, error) {
	requests, err := r.client.ConnectionRequest.Query().
		Where(
			connectionrequest.And(
				connectionrequest.ReceiverIDEQ(userID),
				connectionrequest.StatusEQ(connectionrequest.StatusPending),
			),
		).
		WithSender().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending requests: %w", err)
	}
	return requests, nil
}

func (r *connectionRequestRepository) GetSentRequests(ctx context.Context, userID uuid.UUID) ([]*ent.ConnectionRequest, error) {
	requests, err := r.client.ConnectionRequest.Query().
		Where(connectionrequest.SenderIDEQ(userID)).
		WithReceiver().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get sent requests: %w", err)
	}
	return requests, nil
}

func (r *connectionRequestRepository) GetReceivedRequests(ctx context.Context, userID uuid.UUID) ([]*ent.ConnectionRequest, error) {
	requests, err := r.client.ConnectionRequest.Query().
		Where(connectionrequest.ReceiverIDEQ(userID)).
		WithSender().
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get received requests: %w", err)
	}
	return requests, nil
}

func (r *connectionRequestRepository) AcceptRequest(ctx context.Context, requestID uuid.UUID) (*ent.ConnectionRequest, *ent.Connection, error) {
	// Start a transaction
	tx, err := r.client.Tx(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to start transaction: %w", err)
	}

	// Get the request
	request, err := tx.ConnectionRequest.Get(ctx, requestID)
	if err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to get connection request: %w", err)
	}

	// Update request status
	updatedRequest, err := tx.ConnectionRequest.UpdateOneID(requestID).
		SetStatus(connectionrequest.StatusAccepted).
		SetRespondedAt(time.Now()).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to update request status: %w", err)
	}

	// Create connection
	newConnection, err := tx.Connection.Create().
		SetUserAID(request.SenderID).
		SetUserBID(request.ReceiverID).
		SetStatus(connection.StatusConnected).
		Save(ctx)
	if err != nil {
		tx.Rollback()
		return nil, nil, fmt.Errorf("failed to create connection: %w", err)
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return nil, nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return updatedRequest, newConnection, nil
}

func (r *connectionRequestRepository) DeclineRequest(ctx context.Context, requestID uuid.UUID) (*ent.ConnectionRequest, error) {
	request, err := r.client.ConnectionRequest.UpdateOneID(requestID).
		SetStatus(connectionrequest.StatusDeclined).
		SetRespondedAt(time.Now()).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to decline request: %w", err)
	}
	return request, nil
}

func (r *connectionRequestRepository) ExpireOldRequests(ctx context.Context) (int, error) {
	// Expire requests older than 30 days that are still pending
	expireThreshold := time.Now().AddDate(0, 0, -30)

	count, err := r.client.ConnectionRequest.Update().
		Where(
			connectionrequest.And(
				connectionrequest.StatusEQ(connectionrequest.StatusPending),
				connectionrequest.CreatedAtLT(expireThreshold),
			),
		).
		SetStatus(connectionrequest.StatusExpired).
		Save(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to expire old requests: %w", err)
	}
	return count, nil
}