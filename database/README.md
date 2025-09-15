# database

**Database** is an event-sourcing database written in Go that stores events conforming to the [CNCF CloudEvents specification](https://cloudevents.io). It provides an in-memory event store with indexing capabilities and JSON-based persistence for reliable event storage and retrieval.

---

## Features

- **In-memory event storage** with fast retrieval by ID, type, and subject
- **Event indexing** for efficient querying by type and subject
- **JSON persistence** to disk for data durability
- **CloudEvents-compatible** event format
- **Graceful shutdown** with automatic data persistence
- **Docker-ready** for easy deployment
- **Event sourcing patterns** with chronological event ordering

---

## Quick Start

### Run with Docker

```bash
docker run \
	-e DATA_DIR=/data \
	-p 5000:5000 \
	-v $(pwd)/data:/data \
	--name eventdb github.com/nicograef/cloudevents/database
```

### Run Locally

```sh
go run .
```

The server will start on port 5000 by default. You can customize the port and data directory using environment variables:

```sh
PORT=8080 go run .
```

```sh
DATA_DIR=/path/to/data go run .
```

---

## Configuration

You can configure the database using environment variables:

| Variable   | Default | Description                    |
|------------|---------|--------------------------------|
| `PORT`     | `5000`  | Port for HTTP server           |
| `DATA_DIR` | `.`     | Directory for data persistence |

---

## API Documentation

The database provides both a Go API and HTTP API for event storage and retrieval.

### HTTP API

#### Add Event

**POST /add**

**Content-Type:** `application/json`

**Payload Example:**

```json
{
  "type": "com.example.user.created:v1",
  "source": "https://api.example.com",
  "subject": "/users/12345",
  "data": {
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

**Success Response:**

```json
{
  "ok": true,
  "event": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "com.example.user.created:v1",
    "time": "2025-09-14T12:34:56Z",
    "source": "https://api.example.com",
    "subject": "/users/12345",
    "data": {
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

**Error Response:**

```json
{
  "ok": false,
  "error": "event type must be at least 5 characters long"
}
```

### Go API

#### Add Event

```go
import "github.com/nicograef/cloudevents/event"

// Create an event candidate
candidate := event.Candidate{
    Type:    "com.example.user.created:v1",
    Source:  "https://api.example.com",
    Subject: "/users/12345",
    Data:    map[string]interface{}{"name": "John Doe", "email": "john@example.com"},
}

// Add to database
event, err := db.AddEvent(candidate)
```

#### Retrieve Events

```go
// Get event by ID
event := db.GetEvent(eventID)

// Get all events (sorted by timestamp)
allEvents := db.GetEvents()

// Get events by type
userEvents := db.GetEventsByType("com.example.user.created:v1")

// Get events by subject
subjectEvents := db.GetEventsBySubject("/users/12345")
```

## Event Format

Events follow the CloudEvents specification:

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "type": "com.example.user.created:v1",
  "time": "2025-09-14T12:34:56Z",
  "source": "https://api.example.com",
  "subject": "/users/12345",
  "data": {
    "name": "John Doe",
    "email": "john@example.com"
  }
}
```

---

## Example: Add Event with curl

```sh
curl -X POST http://localhost:5000/add \
    -H "Content-Type: application/json" \
    -d '{
        "type": "com.example.user.created:v1",
        "source": "https://api.example.com",
        "subject": "/users/12345",
        "data": {
            "name": "John Doe",
            "email": "john@example.com"
        }
    }'
```

**Response:**

```json
{
  "ok": true,
  "event": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "com.example.user.created:v1",
    "time": "2025-09-14T12:34:56Z",
    "source": "https://api.example.com",
    "subject": "/users/12345",
    "data": {
      "name": "John Doe",
      "email": "john@example.com"
    }
  }
}
```

---

## Development

### Install

```bash
# From another module
go get github.com/nicograef/cloudevents/database
```

If developing locally in a monorepo, add a replace directive in your `go.mod`:

```go
replace github.com/nicograef/cloudevents/database => ../database
```

### Build & Run

```sh
go build -o database .
./database
```

### Build Docker Image

```sh
docker build -t github.com/nicograef/cloudevents/database .
```

### Run Docker Container

```sh
docker run -p 5000:5000 -e DATA_DIR=/data -v $(pwd)/data:/data --name database github.com/nicograef/cloudevents/database
```

### Run Tests

```sh
go test ./...
```

---

## Architecture

The database consists of:

- **Events Map**: Primary storage indexed by event ID
- **Type Index**: Secondary index for fast type-based queries
- **Subject Index**: Secondary index for fast subject-based queries
- **Persistence Layer**: JSON serialization to/from disk

The persistence format stores events as a JSON array for efficient parsing and minimal overhead.

---

## Related Projects

- [CloudEvents spec](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md)
- [Event module](../event) - CloudEvents-compatible event library
- [Queue module](../queue) - Message queue for event delivery
- [EventStore](https://www.eventstore.com/) - Production-grade event sourcing database

---

## License

MIT
