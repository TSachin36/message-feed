package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"time"

	"message-feed/internal/models"
	"message-feed/internal/storage"
)

const dataFile = "data/messages.txt"

func main() {

	user := flag.String("user", "", "User ID")
	message := flag.String("message", "", "Message")

	flag.Parse()

	logger := slog.Default()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	traceID := fmt.Sprintf(
		"TRACE-%d",
		time.Now().UnixNano(),
	)

	ctx := context.WithValue(
		context.Background(),
		"traceID",
		traceID,
	)

	if *user == "" || *message == "" {
		fmt.Println("Please provide both -user and -message")
		return
	}

	msg := models.Message{
		UserID: *user,
		Text:   *message,
	}

	err := storage.SaveMessage(ctx, dataFile, msg)
	if err != nil {
		logger.Error("Failed to save message", "error", err)
		return
	}

	messages, err := storage.GetLast10Messages(ctx, dataFile)
	if err != nil {
		logger.Error(
			"Failed to read messages",
			"error", err,
		)
		return
	}

	fmt.Println("\nLast Messages:")

	for _, message := range messages {
		fmt.Printf(
			"%s: %s\n",
			message.UserID,
			message.Text,
		)
	}
	fmt.Println("\nApplication is running...")
	fmt.Println("Press Ctrl+C to exit.")

	<-sigChan

	logger.Info("Interrupt signal received. Shutting down...")
}
