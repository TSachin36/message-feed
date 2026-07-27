package api

import (
	"context"

	"message-feed/internal/models"
	pb "message-feed/proto"
)

func getMessages(
	ctx context.Context,
	userID string,
) ([]models.Message, error) {

	response, err := grpcClient.GetLast10(
		ctx,
		&pb.UserRequest{
			UserID: userID,
		},
	)
	if err != nil {
		return nil, err
	}

	return toModels(response.Messages), nil
}

func saveMessage(
	ctx context.Context,
	msg models.Message,
) error {

	_, err := grpcClient.Save(
		ctx,
		&pb.Message{
			UserID: msg.UserID,
			Text:   msg.Text,
		},
	)

	return err
}
