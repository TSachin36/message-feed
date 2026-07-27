package main

import (
	"net/http"
)

func websocketHandler(w http.ResponseWriter, r *http.Request) {

	userID := r.URL.Query().Get("user")

	if userID == "" {
		http.Error(
			w,
			"Missing user",
			http.StatusBadRequest,
		)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.Error(
			"Failed to upgrade connection",
			"error", err,
		)
		return
	}

	logger.Info(
		"WebSocket client connected",
		"user", userID,
	)

	messages, err := getMessages(
		r.Context(),
		userID,
	)
	if err != nil {

		logger.Error(
			"Failed to read messages",
			"error", err,
		)

		conn.Close()

		return
	}

	err = conn.WriteJSON(messages)
	if err != nil {
		logger.Error(
			"Failed to send messages",
			"error", err,
		)

		conn.Close()

		return
	}

	logger.Info(
		"Sent last 10 messages",
		"user", userID,
	)

	register <- client{
		conn:   conn,
		userID: userID,
	}

	for {
		if _, _, err := conn.ReadMessage(); err != nil {

			unregister <- conn

			break
		}
	}
}
