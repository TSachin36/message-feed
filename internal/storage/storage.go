package storage

import (
	"context"
	"fmt"
	"log/slog"
	"message-feed/internal/models"
	"os"
	"strings"
)

func SaveMessage(ctx context.Context, filename string, msg models.Message) error {

	traceID, _ := ctx.Value("traceID").(string)

	logger := slog.Default()

	file, err := os.OpenFile(
		filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)

	if err != nil {
		return err
	}

	defer file.Close()

	_, err = file.WriteString(
		fmt.Sprintf("%s: %s\n", msg.UserID, msg.Text),
	)

	if err != nil {
		return err
	}

	logger.Info(
		"Message Saved",
		"traceID", traceID,
		"user", msg.UserID,
		"message", msg.Text,
	)

	return nil
}

func GetLast10Messages(ctx context.Context, filename string) ([]models.Message, error) {

	traceID, _ := ctx.Value("traceID").(string)

	logger := slog.Default()

	logger.Info(
		"Reading last 10 messages",
		"traceID", traceID,
	)

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	text := string(data)

	lines := strings.Split(text, "\n")

	if len(lines) > 0 && lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}

	var messages []models.Message

	for _, line := range lines {

		parts := strings.SplitN(line, ":", 2)

		if len(parts) != 2 {
			continue
		}

		msg := models.Message{
			UserID: strings.TrimSpace(parts[0]),
			Text:   strings.TrimSpace(parts[1]),
		}

		messages = append(messages, msg)
	}

	if len(messages) <= 10 {
		return messages, nil
	}

	return messages[len(messages)-10:], nil
}

func GetLast10MessagesForUser(
	ctx context.Context,
	filename string,
	userID string,
) ([]models.Message, error) {

	traceID, _ := ctx.Value("traceID").(string)

	logger := slog.Default()

	logger.Info(
		"Reading last 10 messages for user",
		"traceID", traceID,
		"user", userID,
	)

	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	text := string(data)
	lines := strings.Split(text, "\n")

	var messages []models.Message

	for _, line := range lines {

		if line == "" {
			continue
		}

		parts := strings.SplitN(line, ":", 2)

		if len(parts) != 2 {
			continue
		}

		msg := models.Message{
			UserID: strings.TrimSpace(parts[0]),
			Text:   strings.TrimSpace(parts[1]),
		}

		if msg.UserID != userID {
			continue
		}

		messages = append(messages, msg)
	}

	if len(messages) <= 10 {
		return messages, nil
	}

	return messages[len(messages)-10:], nil
}
