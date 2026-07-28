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
var dataFile string

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
		dataFile,
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
		dataFile,
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

func Run(
	address string,
	filename string,
) {

	dataFile = filename

	lis, err := net.Listen(
		"tcp",
		address,
	)
	if err != nil {
		log.Printf(
			"Failed to start gRPC store on %s: %v",
			address,
			err,
		)
		return
	}

	grpcServer = grpc.NewServer()

	pb.RegisterMessageStoreServer(
		grpcServer,
		&server{},
	)

	log.Printf(
		"gRPC Store Server running on %s using %s",
		address,
		filename,
	)

	if err := grpcServer.Serve(lis); err != nil {
		log.Printf(
			"gRPC store stopped: %v",
			err,
		)
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
