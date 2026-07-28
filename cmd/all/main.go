package main

import (
	"os"
	"os/signal"

	"message-feed/internal/api"
	"message-feed/internal/store"
)

func main() {

	store.Run(
		":50051",
		"data/messages.txt",
	)

	api.Start()

	sigChan := make(chan os.Signal, 1)

	signal.Notify(
		sigChan,
		os.Interrupt,
	)

	<-sigChan

	api.Shutdown()
	store.Shutdown()
}
