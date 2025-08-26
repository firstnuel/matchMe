package connections

import (
	"context"
	"fmt"
	"match-me/ent"
	"match-me/ent/connection"

	"github.com/google/uuid"
)

type connectionRepository struct {
	client *ent.Client
}

func NewConnectionRepository(client *ent.Client) ConnectionRepository {
	return &connectionRepository{
		client: client,
	}
}

func (r *connectionRepository) CreateConnection(ctx context.Context, userAID, userBID uuid.UUID) (*ent.Connection, error) {
	conn, err := r.client.Connection.Create().
		SetUserAID(userAID).
		SetUserBID(userBID).
		SetStatus(connection.StatusConnected).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection: %w", err)
	}
	return conn, nil
}

func (r *connectionRepository) GetConnection(ctx context.Context, connectionID uuid.UUID) (*ent.Connection, error) {
	conn, err := r.client.Connection.Get(ctx, connectionID)
	if err != nil {
		return nil, fmt.Errorf("failed to get connection: %w", err)
	}
	return conn, nil
}

func (r *connectionRepository) GetConnectionBetweenUsers(ctx context.Context, userAID, userBID uuid.UUID) (*ent.Connection, error) {
	conn, err := r.client.Connection.Query().
		Where(
			connection.Or(
				connection.And(
					connection.UserAIDEQ(userAID),
					connection.UserBIDEQ(userBID),
				),
				connection.And(
					connection.UserAIDEQ(userBID),
					connection.UserBIDEQ(userAID),
				),
			),
		).
		First(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get connection between users: %w", err)
	}
	return conn, nil
}

func (r *connectionRepository) UpdateConnectionStatus(ctx context.Context, connectionID uuid.UUID, status string) (*ent.Connection, error) {
	conn, err := r.client.Connection.UpdateOneID(connectionID).
		SetStatus(connection.Status(status)).
		Save(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to update connection status: %w", err)
	}
	return conn, nil
}

func (r *connectionRepository) DeleteConnection(ctx context.Context, connectionID uuid.UUID) error {
	err := r.client.Connection.DeleteOneID(connectionID).Exec(ctx)
	if err != nil {
		return fmt.Errorf("failed to delete connection: %w", err)
	}
	return nil
}

func (r *connectionRepository) GetUserConnections(ctx context.Context, userID uuid.UUID) ([]*ent.Connection, error) {
	connections, err := r.client.Connection.Query().
		Where(
			connection.Or(
				connection.UserAIDEQ(userID),
				connection.UserBIDEQ(userID),
			),
		).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get user connections: %w", err)
	}
	return connections, nil
}

func (r *connectionRepository) GetActiveConnections(ctx context.Context, userID uuid.UUID) ([]*ent.Connection, error) {
	connections, err := r.client.Connection.Query().
		Where(
			connection.And(
				connection.StatusEQ(connection.StatusConnected),
				connection.Or(
					connection.UserAIDEQ(userID),
					connection.UserBIDEQ(userID),
				),
			),
		).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get active connections: %w", err)
	}
	return connections, nil
}

func (r *connectionRepository) GetConnectionsWithUsers(ctx context.Context, userID uuid.UUID) ([]*ent.Connection, error) {
	connections, err := r.client.Connection.Query().
		Where(
			connection.Or(
				connection.UserAIDEQ(userID),
				connection.UserBIDEQ(userID),
			),
		).
		WithUserA(func(uq *ent.UserQuery) {
			uq.WithPhotos()
		}).
		WithUserB(func(uq *ent.UserQuery) {
			uq.WithPhotos()
		}).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connections with users: %w", err)
	}
	return connections, nil
}

func (r *connectionRepository) GetConnectedUserIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	connections, err := r.client.Connection.Query().
		Where(
			connection.And(
				connection.StatusEQ(connection.StatusConnected),
				connection.Or(
					connection.UserAIDEQ(userID),
					connection.UserBIDEQ(userID),
				),
			),
		).
		All(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get connected user IDs: %w", err)
	}

	// Extract the other user IDs from each connection
	var connectedUserIDs []uuid.UUID
	for _, conn := range connections {
		if conn.UserAID == userID {
			connectedUserIDs = append(connectedUserIDs, conn.UserBID)
		} else {
			connectedUserIDs = append(connectedUserIDs, conn.UserAID)
		}
	}

	return connectedUserIDs, nil
}
