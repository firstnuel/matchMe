package connections

import (
	"context"
	"fmt"
	"match-me/ent"
	"match-me/internal/models"
	"match-me/internal/repositories/connections"

	"github.com/google/uuid"
)

type connectionRequestUsecase struct {
	requestRepo    connections.ConnectionRequestRepository
	connectionRepo connections.ConnectionRepository
}

func NewConnectionRequestUsecase(
	requestRepo connections.ConnectionRequestRepository,
	connectionRepo connections.ConnectionRepository,
) ConnectionRequestUsecase {
	return &connectionRequestUsecase{
		requestRepo:    requestRepo,
		connectionRepo: connectionRepo,
	}
}

func (u *connectionRequestUsecase) SendRequest(ctx context.Context, senderID, receiverID uuid.UUID, message string) (*models.ConnectionRequest, error) {
	// Validate that users are not the same
	if senderID == receiverID {
		return nil, fmt.Errorf("cannot send connection request to yourself")
	}

	// Check if connection already exists
	existingConnection, err := u.connectionRepo.GetConnectionBetweenUsers(ctx, senderID, receiverID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing connection: %w", err)
	}
	if existingConnection != nil {
		return nil, fmt.Errorf("connection already exists between users")
	}

	// Check if request already exists
	existingRequest, err := u.requestRepo.GetConnectionRequestBetweenUsers(ctx, senderID, receiverID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing request: %w", err)
	}
	if existingRequest != nil {
		return nil, fmt.Errorf("connection request already exists")
	}

	// Create the request
	entRequest, err := u.requestRepo.CreateConnectionRequest(ctx, senderID, receiverID, message)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection request: %w", err)
	}

	return models.ToConnectionRequest(entRequest), nil
}

func (u *connectionRequestUsecase) GetPendingRequests(ctx context.Context, userID uuid.UUID) ([]*models.ConnectionRequest, error) {
	entRequests, err := u.requestRepo.GetPendingRequestsForUser(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get pending requests: %w", err)
	}

	return models.ToConnectionRequests(entRequests), nil
}

func (u *connectionRequestUsecase) AcceptRequest(ctx context.Context, userID, requestID uuid.UUID) (*models.Connection, error) {
	// Get the request to verify ownership
	request, err := u.requestRepo.GetConnectionRequest(ctx, requestID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("connection request not found")
		}
		return nil, fmt.Errorf("failed to get connection request: %w", err)
	}

	// Verify that the user is the receiver of this request
	if request.ReceiverID != userID {
		return nil, fmt.Errorf("unauthorized: user is not the receiver of this request")
	}

	// Verify request is still pending
	if request.Status != "pending" {
		return nil, fmt.Errorf("request is no longer pending")
	}

	// Accept the request (this creates the connection)
	_, entConnection, err := u.requestRepo.AcceptRequest(ctx, requestID)
	if err != nil {
		return nil, fmt.Errorf("failed to accept request: %w", err)
	}

	return models.ToConnection(entConnection), nil
}

func (u *connectionRequestUsecase) DeclineRequest(ctx context.Context, userID, requestID uuid.UUID) error {
	// Get the request to verify ownership
	request, err := u.requestRepo.GetConnectionRequest(ctx, requestID)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("connection request not found")
		}
		return fmt.Errorf("failed to get connection request: %w", err)
	}

	// Verify that the user is the receiver of this request
	if request.ReceiverID != userID {
		return fmt.Errorf("unauthorized: user is not the receiver of this request")
	}

	// Verify request is still pending
	if request.Status != "pending" {
		return fmt.Errorf("request is no longer pending")
	}

	// Decline the request
	_, err = u.requestRepo.DeclineRequest(ctx, requestID)
	if err != nil {
		return fmt.Errorf("failed to decline request: %w", err)
	}

	return nil
}