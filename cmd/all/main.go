package main

import (
	"os"
	"os/signal"

	"message-feed/internal/api"
	"message-feed/internal/store"
)

func main() {

	go store.Run()

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
