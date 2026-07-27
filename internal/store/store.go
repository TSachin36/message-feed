package store

import (
	"context"
	"log"
	"net"

	"message-feed/internal/models"
	"message-feed/internal/storage"
	pb "message-feed/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedMessageStoreServer
}

var grpcServer *grpc.Server

// gRPC requires this receiver method.
// It delegates the real work to the package-level save function.
func (s *server) Save(
	ctx context.Context,
	req *pb.Message,
) (*pb.Empty, error) {
	return save(ctx, req)
}

// gRPC requires this receiver method.
// It delegates the real work to the package-level getLast10 function.
func (s *server) GetLast10(
	ctx context.Context,
	req *pb.UserRequest,
) (*pb.MessageList, error) {
	return getLast10(ctx, req)
}

func save(
	ctx context.Context,
	req *pb.Message,
) (*pb.Empty, error) {

	msg := models.Message{
		UserID: req.UserID,
		Text:   req.Text,
	}

	err := storage.SaveMessage(
		ctx,
		"data/messages.txt",
		msg,
	)
	if err != nil {
		return nil, err
	}

	return &pb.Empty{}, nil
}

func getLast10(
	ctx context.Context,
	req *pb.UserRequest,
) (*pb.MessageList, error) {

	messages, err := storage.GetLast10MessagesForUser(
		ctx,
		"data/messages.txt",
		req.UserID,
	)
	if err != nil {
		return nil, err
	}

	var pbMessages []*pb.Message

	for _, msg := range messages {
		pbMessages = append(
			pbMessages,
			&pb.Message{
				UserID: msg.UserID,
				Text:   msg.Text,
			},
		)
	}

	return &pb.MessageList{
		Messages: pbMessages,
	}, nil
}

func Run() {

	lis, err := net.Listen(
		"tcp",
		":50051",
	)
	if err != nil {
		log.Printf("Failed to start gRPC store: %v", err)
		return
	}

	grpcServer = grpc.NewServer()

	pb.RegisterMessageStoreServer(
		grpcServer,
		&server{},
	)

	log.Println("gRPC Store Server running on :50051")

	if err := grpcServer.Serve(lis); err != nil {
		log.Printf("gRPC store stopped: %v", err)
	}
}

func Shutdown() {

	if grpcServer == nil {
		return
	}

	log.Println("Shutting down gRPC Store Server...")

	grpcServer.GracefulStop()

	log.Println("gRPC Store Server stopped gracefully")
}
