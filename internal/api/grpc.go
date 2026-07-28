package api

import (
	"log"

	pb "message-feed/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var grpcConnections = make(
	map[string]*grpc.ClientConn,
)

var grpcClients = make(
	map[string]pb.MessageStoreClient,
)

func connectGRPC() {

	buildHashRing()

	for _, address := range shardAddresses {

		conn, err := grpc.Dial(
			address,
			grpc.WithTransportCredentials(
				insecure.NewCredentials(),
			),
		)
		if err != nil {
			log.Fatal(err)
		}

		grpcConnections[address] = conn

		grpcClients[address] =
			pb.NewMessageStoreClient(conn)

		logger.Info(
			"Connected to shard",
			"address", address,
		)
	}
}

func closeGRPC() {

	for address, conn := range grpcConnections {

		if err := conn.Close(); err != nil {

			logger.Error(
				"Failed to close shard connection",
				"address", address,
				"error", err,
			)
		}
	}
}