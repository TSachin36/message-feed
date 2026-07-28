# Message Feed

A Go message-feed application demonstrating file storage, HTTP APIs, WebSockets, gRPC, structured logging, concurrency, benchmarking, profiling, and consistent-hash sharding.

## Features

- CLI and interactive CRUD REPL
- Multi-user message storage
- `log/slog` structured logging and Trace IDs
- HTTP API using `net/http`
- Static pages and HTML templates
- WebSocket live subscriptions
- gRPC storage service
- Actor/CSP concurrency
- Graceful shutdown
- Unit, parallel, and benchmark tests
- HTTP load testing
- pprof profiling
- Consistent hashing with virtual nodes
- Multiple gRPC storage backends

## Architecture

```text
                    Clients
               HTTP / WebSocket
                      |
                      v
                API Server :8080
                      |
                hash(userID)
                      |
              Consistent Hash Ring
                100 virtual nodes
                      |
          +-----------+-----------+
          |           |           |
          v           v           v
      gRPC :50051 gRPC :50052 gRPC :50053
          |           |           |
          v           v           v
      shard1.txt  shard2.txt  shard3.txt
```

The frontend routes each user to a storage backend using consistent hashing. Messages are persisted through gRPC and then distributed to live subscribers using the Actor/CSP pattern.

## Project Structure

```text
cmd/
  all/        Combined startup
  api/        HTTP/WebSocket frontend
  cli/        CLI
  client/     WebSocket client
  loadtest/   Load tester
  repl/       CRUD REPL
  store/      gRPC storage backend

internal/
  api/        HTTP, WebSocket and sharding logic
  models/     Data models
  storage/    File storage
  store/      gRPC service

proto/        Protocol Buffer definitions/generated code
templates/    HTML templates
web/          Static content
data/         Runtime data
```

## CLI

Save a message and display the last 10 messages:

```bash
go run ./cmd/cli -user="alice" -message="Hello from CLI"
```

The CLI remains running until `Ctrl+C`.

## REPL

```bash
go run ./cmd/repl
```

Commands:

```text
create <user> <message>
list
update <number> <user> <message>
delete <number>
help
exit
```

## Running the Sharded Application

Start three storage backends in separate terminals:

```bash
go run ./cmd/store "-port=50051" "-data=data/shard1.txt"
go run ./cmd/store "-port=50052" "-data=data/shard2.txt"
go run ./cmd/store "-port=50053" "-data=data/shard3.txt"
```

Then start the frontend:

```bash
go run ./cmd/api
```

The API runs on port `8080` and pprof on port `6060`.

## HTTP API

Create a message:

```text
POST /messages
```

PowerShell example:

```powershell
Invoke-RestMethod `
    -Method Post `
    -Uri "http://localhost:8080/messages" `
    -ContentType "application/json" `
    -Body '{"user":"alice","text":"Hello"}'
```

Retrieve a user's messages:

```text
GET /messages?user=alice
```

```powershell
Invoke-RestMethod `
    -Method Get `
    -Uri "http://localhost:8080/messages?user=alice"
```

## Web Pages

Static page:

```text
http://localhost:8080/about/
```

Dynamic user message page:

```text
http://localhost:8080/list?user=alice
```

The `user` parameter is required for `/list`.

## WebSockets

The WebSocket endpoint provides live message updates to subscribers.

```text
ws://localhost:8080/ws?user=alice
```

Run the client with:

```bash
go run ./cmd/client
```

The server sends recent messages and keeps the connection open for new messages.

## gRPC

The storage service defines:

```text
Save
GetLast10
```

Storage runs independently from the HTTP frontend, allowing multiple backend instances.

## Actor / CSP Concurrency

WebSocket clients are coordinated using channels:

```text
register
unregister
broadcast
```

A dedicated actor goroutine owns the client state, avoiding concurrent mutation of the client map.

Messages are broadcast after they are successfully persisted.

## Consistent Hash Sharding

Users are routed to storage backends using CRC32 consistent hashing.

Each physical backend has **100 virtual nodes**, giving three backends 300 positions on the hash ring.

A 10,000-user distribution test produced:

```text
localhost:50051 -> 4070
localhost:50052 -> 3205
localhost:50053 -> 2725
```

Adding a fourth shard moved:

```text
2615 / 10000 users (26.15%)
```

Removing `localhost:50052` moved exactly the 3205 users previously assigned to it, while users on unaffected shards remained in place.

The ring handles routing only; automatic migration of historical data between shards is not implemented.

## Logging and Trace IDs

The application uses `log/slog` for structured logging and `context.Context` for Trace ID propagation.

Example:

```text
INFO Message Saved traceID=TEST-456 user=alice message="Hello"
```

## Graceful Shutdown

Interrupt signals are handled using `os/signal`.

The HTTP server uses `http.Server.Shutdown` to stop accepting new requests while allowing active requests to finish.

The gRPC store also supports graceful shutdown.

## Published Core Module

The reusable Actor/CSP logic is published as:

```text
github.com/TSachin36/message-feed-core
```

Version:

```text
v1.0.0
```

## Testing

Run all tests:

```bash
go test ./...
```

The project includes unit tests, parallel concurrency tests, and consistent-hashing tests.

Final validation:

```bash
go test ./...
go vet ./...
go build ./...
```

## Benchmarks

Run storage benchmarks:

```bash
go test ./internal/storage -bench=. -benchmem
```

Benchmarks cover message saving and retrieving the last 10 messages for a user.

## Load Testing

Run the HTTP load tester:

```bash
go run ./cmd/loadtest
```

Example result:

```text
Requests:     1000
Concurrency:  20
Successful:   1000
Failed:       0
Requests/sec: 2734.77
```

## pprof

Profiling is available at:

```text
http://localhost:6060/debug/pprof/
```

CPU profile:

```bash
go tool pprof "http://localhost:6060/debug/pprof/profile?seconds=15"
```

Allocation profile:

```bash
go tool pprof "http://localhost:6060/debug/pprof/allocs"
```

## Runtime Data

Shard files are generated at runtime and ignored by Git:

```gitignore
data/shard*.txt
```

## Assignment Coverage

The project covers the required progression:

- CLI and file storage
- Structured logging and Trace IDs
- Unit and parallel tests
- Signal handling and graceful shutdown
- HTTP API and middleware
- Static and dynamic web pages
- WebSockets and client
- Protocol Buffers and gRPC
- Actor/CSP concurrency
- CRUD REPL
- Multiple startup modes
- Multi-user support
- Published reusable module
- Benchmarks and load testing
- pprof profiling
- Frontend/backend separation
- Consistent-hash sharding with virtual nodes

## Technologies

Go, `net/http`, `html/template`, `log/slog`, `context`, goroutines, channels, gRPC, Protocol Buffers, Gorilla WebSocket, pprof, and consistent hashing.

## Author

Sachin