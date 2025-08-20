package connections

import (
	"context"
	"fmt"
	"match-me/ent"
	"match-me/internal/models"
	"match-me/internal/repositories/connections"

	"github.com/google/uuid"
)

type connectionUsecase struct {
	connectionRepo connections.ConnectionRepository
}

func NewConnectionUsecase(connectionRepo connections.ConnectionRepository) ConnectionUsecase {
	return &connectionUsecase{
		connectionRepo: connectionRepo,
	}
}

func (u *connectionUsecase) GetUserConnections(ctx context.Context, userID uuid.UUID) ([]*models.Connection, error) {
	// Get connections with user details
	entConnections, err := u.connectionRepo.GetConnectionsWithUsers(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user connections: %w", err)
	}

	// Convert to models
	return models.ToConnections(entConnections), nil
}

func (u *connectionUsecase) DeleteConnection(ctx context.Context, userID, connectionID uuid.UUID) error {
	// Get the connection to verify ownership
	connection, err := u.connectionRepo.GetConnection(ctx, connectionID)
	if err != nil {
		if ent.IsNotFound(err) {
			return fmt.Errorf("connection not found")
		}
		return fmt.Errorf("failed to get connection: %w", err)
	}

	// Verify that the user is part of this connection
	if connection.UserAID != userID && connection.UserBID != userID {
		return fmt.Errorf("unauthorized: user is not part of this connection")
	}

	// Delete the connection
	err = u.connectionRepo.DeleteConnection(ctx, connectionID)
	if err != nil {
		return fmt.Errorf("failed to delete connection: %w", err)
	}

	return nil
}