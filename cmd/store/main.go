package main

import (
	"flag"
	"fmt"

	"message-feed/internal/store"
)

func main() {

	port := flag.Int(
		"port",
		50051,
		"gRPC store port",
	)

	data := flag.String(
		"data",
		"data/messages.txt",
		"message storage file",
	)

	flag.Parse()

	address := fmt.Sprintf(
		":%d",
		*port,
	)

	store.Run(
		address,
		*data,
	)
}
