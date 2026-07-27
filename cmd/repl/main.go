package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"message-feed/internal/models"
)

func main() {
	scanner := bufio.NewScanner(os.Stdin)

	var messages []models.Message

	fmt.Println("Message Feed REPL")
	fmt.Println("Commands: create, list, update, delete, help, exit")

	for {
		fmt.Print("> ")

		if !scanner.Scan() {
			break
		}

		input := strings.TrimSpace(scanner.Text())

		if input == "" {
			continue
		}

		parts := strings.SplitN(input, " ", 2)
		command := strings.ToLower(parts[0])

		switch command {
		case "create":
			if len(parts) != 2 {
				fmt.Println("Usage: create <user> <message>")
				continue
			}

			createParts := strings.SplitN(parts[1], " ", 2)
			if len(createParts) != 2 {
				fmt.Println("Usage: create <user> <message>")
				continue
			}

			msg := models.Message{
				UserID: createParts[0],
				Text:   createParts[1],
			}

			messages = append(messages, msg)

			fmt.Printf("Created message %d\n", len(messages))

		case "list":
			if len(messages) == 0 {
				fmt.Println("No messages")
				continue
			}

			for index, msg := range messages {
				fmt.Printf(
					"%d. %s: %s\n",
					index+1,
					msg.UserID,
					msg.Text,
				)
			}

		case "update":
			if len(parts) != 2 {
				fmt.Println("Usage: update <number> <user> <message>")
				continue
			}

			updateParts := strings.SplitN(parts[1], " ", 3)
			if len(updateParts) != 3 {
				fmt.Println("Usage: update <number> <user> <message>")
				continue
			}

			number, err := strconv.Atoi(updateParts[0])
			if err != nil || number < 1 || number > len(messages) {
				fmt.Println("Invalid message number")
				continue
			}

			messages[number-1] = models.Message{
				UserID: updateParts[1],
				Text:   updateParts[2],
			}

			fmt.Printf("Updated message %d\n", number)

		case "delete":
			if len(parts) != 2 {
				fmt.Println("Usage: delete <number>")
				continue
			}

			number, err := strconv.Atoi(parts[1])
			if err != nil || number < 1 || number > len(messages) {
				fmt.Println("Invalid message number")
				continue
			}

			index := number - 1

			messages = append(
				messages[:index],
				messages[index+1:]...,
			)

			fmt.Printf("Deleted message %d\n", number)

		case "help":
			fmt.Println("create <user> <message>")
			fmt.Println("list")
			fmt.Println("update <number> <user> <message>")
			fmt.Println("delete <number>")
			fmt.Println("help")
			fmt.Println("exit")

		case "exit", "quit":
			fmt.Println("Goodbye")
			return

		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read input: %v\n", err)
	}
}
