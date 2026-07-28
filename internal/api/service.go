package api

import (
	"context"
	"fmt"

	"message-feed/internal/models"
	pb "message-feed/proto"
)

func getMessages(
	ctx context.Context,
	userID string,
) ([]models.Message, error) {

	address := shardForUser(userID)

	client, exists := grpcClients[address]
	if !exists {
		return nil, fmt.Errorf(
			"no gRPC client for shard %s",
			address,
		)
	}

	response, err := client.GetLast10(
		ctx,
		&pb.UserRequest{
			UserID: userID,
		},
	)
	if err != nil {
		return nil, err
	}

	logger.Info(
		"Read routed to shard",
		"user", userID,
		"shard", address,
	)

	return toModels(response.Messages), nil
}

func saveMessage(
	ctx context.Context,
	msg models.Message,
) error {

	address := shardForUser(msg.UserID)

	client, exists := grpcClients[address]
	if !exists {
		return fmt.Errorf(
			"no gRPC client for shard %s",
			address,
		)
	}

	_, err := client.Save(
		ctx,
		&pb.Message{
			UserID: msg.UserID,
			Text:   msg.Text,
		},
	)
	if err != nil {
		return err
	}

	logger.Info(
		"Write routed to shard",
		"user", msg.UserID,
		"shard", address,
	)

	return nil
}
