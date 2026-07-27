# Message Feed

A Go application that demonstrates file storage, HTTP APIs, WebSockets, gRPC, structured logging, middleware, templates, concurrency using the Actor/CSP pattern, and graceful shutdown.

---

## Project Structure

```
message-feed/
│
├── cmd/
│   ├── cli/          # Command-line application
│   ├── server/       # HTTP + WebSocket server
│   ├── store/        # gRPC storage service
│   └── client/       # WebSocket client
│
├── internal/
│   ├── models/
│   └── storage/
│
├── proto/
│   ├── message.proto
│   ├── message.pb.go
│   └── message_grpc.pb.go
│
├── templates/
│   └── messages.html
│
├── web/
│   └── about/
│
└── data/
    └── messages.txt
```

---

## Features

- CLI application
- File-based message storage
- Structured logging using `log/slog`
- Context propagation with Trace IDs
- Unit tests
- Parallel tests
- Graceful shutdown
- HTTP REST API
- Middleware
- Static web page
- Dynamic HTML template
- WebSocket support
- WebSocket client
- gRPC storage service
- Actor/CSP pattern for concurrent WebSocket broadcasting

---

## Architecture

```
                HTTP POST
                     │
                     ▼
             HTTP Server (:8080)
                     │
                     │ gRPC
                     ▼
          gRPC Store Server (:50051)
                     │
                     ▼
             data/messages.txt

                     ▲
                     │
          WebSocket Clients (/ws)
```

---

# Requirements

- Go 1.24+
- Protocol Buffers (`protoc`)
- gorilla/websocket
- google.golang.org/grpc

---

# Install Dependencies

```bash
go mod tidy
```

---

# Generate gRPC Code

```bash
protoc --go_out=. --go-grpc_out=. proto/message.proto
```

---

# Running the Application

Open **three terminals**.

## Terminal 1

Start the gRPC Store.

```bash
go run ./cmd/store
```

Expected output:

```
gRPC Store Server running on :50051
```

---

## Terminal 2

Start the HTTP Server.

```bash
go run ./cmd/server
```

Expected output:

```
Server started
address=:8080
```

---

## Terminal 3

Start the WebSocket Client.

```bash
go run ./cmd/client
```

Expected output:

```
Connected to server

Messages received:

Sachin: Hello
Alice: Testing
...
```

---

# Using the HTTP API

## Get Messages

```
GET /messages
```

Example:

```bash
curl http://localhost:8080/messages
```

---

## Save Message

```
POST /messages
```

Example:

```bash
curl -X POST http://localhost:8080/messages \
-H "Content-Type: application/json" \
-d "{\"userID\":\"Sachin\",\"text\":\"Hello from HTTP\"}"
```

Response:

```json
{
    "userID":"Sachin",
    "text":"Hello from HTTP"
}
```

---

# Web Pages

## Static Page

```
http://localhost:8080/about/
```

Served using `http.FileServer`.

---

## Dynamic Page

```
http://localhost:8080/list
```

Displays the latest 10 messages using Go HTML templates.

---

# WebSocket

Endpoint:

```
ws://localhost:8080/ws
```

When a client connects it:

1. Receives the last 10 stored messages.
2. Remains connected.
3. Receives new messages in real time as they are posted.

---

# Running Tests

Run all tests:

```bash
go test ./...
```

Run with verbose output:

```bash
go test -v ./...
```

---

# Technologies Used

- Go
- gRPC
- Protocol Buffers
- Gorilla WebSocket
- net/http
- html/template
- log/slog
- Context
- Channels
- Goroutines

---

# Concurrency Model

The server uses the Actor / CSP pattern.

Three channels coordinate all WebSocket activity:

- `register`
- `unregister`
- `broadcast`

A dedicated goroutine owns the client map, ensuring concurrent safety without mutexes.

```
HTTP POST
     │
     ▼
broadcast channel
     │
     ▼
Actor Goroutine
     │
     ├── Client 1
     ├── Client 2
     ├── Client 3
     └── Client N
```

---

# Project Goals Completed

- CLI application
- File storage
- Structured logging
- Context propagation
- Unit testing
- Graceful shutdown
- HTTP API
- Middleware
- Static web pages
- HTML templates
- WebSockets
- WebSocket client
- gRPC service
- Protocol Buffers
- Actor/CSP concurrency

---

# Author

Sachin