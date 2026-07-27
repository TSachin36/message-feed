package main

import (
	"log/slog"
	"net/http"

	"github.com/gorilla/websocket"
)

var logger = slog.Default()

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}
