package api

import (
	"message-feed/internal/models"

	"github.com/gorilla/websocket"
)

type client struct {
	conn   *websocket.Conn
	userID string
}

var clients = make(map[*websocket.Conn]string)
var register = make(chan client)
var unregister = make(chan *websocket.Conn)
var broadcast = make(chan models.Message)

func actor() {

	for {

		select {

		case client := <-register:

			clients[client.conn] = client.userID

			logger.Info(
				"Client registered",
				"user", client.userID,
			)

		case conn := <-unregister:

			if _, exists := clients[conn]; exists {

				delete(clients, conn)

				conn.Close()

				logger.Info("Client unregistered")
			}

		case msg := <-broadcast:

			for conn, userID := range clients {

				if userID != msg.UserID {
					continue
				}

				if err := conn.WriteJSON(msg); err != nil {

					delete(clients, conn)

					conn.Close()

					logger.Error(
						"Failed to write to websocket",
						"error", err,
					)
				}
			}
		}
	}
}
