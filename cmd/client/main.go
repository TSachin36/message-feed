package main

import (
	"fmt"
	"log"

	"message-feed/internal/models"

	"github.com/gorilla/websocket"
)

func main() {

	conn, _, err := websocket.DefaultDialer.Dial(
		"ws://localhost:8080/ws?user=alice",
		nil,
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	fmt.Println("Connected to server")

	var messages []models.Message

	err = conn.ReadJSON(&messages)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nLast 10 Messages:")

	for _, msg := range messages {
		fmt.Printf(
			"%s: %s\n",
			msg.UserID,
			msg.Text,
		)
	}

	fmt.Println("\nWaiting for new messages...")

	for {

		var msg models.Message

		err = conn.ReadJSON(&msg)
		if err != nil {
			log.Println(
				"Connection closed:",
				err,
			)
			return
		}

		fmt.Printf(
			"New message -> %s: %s\n",
			msg.UserID,
			msg.Text,
		)
	}
}
