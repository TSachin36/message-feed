package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func waitForShutdown(server *http.Server) {

	sigChan := make(chan os.Signal, 1)

	signal.Notify(
		sigChan,
		os.Interrupt,
	)

	<-sigChan

	logger.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {

		logger.Error(
			"Shutdown failed",
			"error", err,
		)
		return
	}

	logger.Info("Server stopped gracefully")
}
