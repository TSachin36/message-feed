package api

import (
	"log"

	pb "message-feed/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var grpcConn *grpc.ClientConn
var grpcClient pb.MessageStoreClient

func connectGRPC() {

	var err error

	grpcConn, err = grpc.Dial(
		"localhost:50051",
		grpc.WithTransportCredentials(
			insecure.NewCredentials(),
		),
	)
	if err != nil {
		log.Fatal(err)
	}

	grpcClient = pb.NewMessageStoreClient(grpcConn)
}
