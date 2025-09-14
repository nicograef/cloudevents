# queue

**Queue** is a message queue for event-driven applications, written in Go. It follows the [CNCF Cloudevents specification](https://cloudevents.io) and is designed to receive event messages via HTTP, queue them, and deliver them to a configured webhook. Qugo uses Go channels to decouple producers and consumers, ensuring reliable delivery and easy integration.

---

## Features

- **Async message queue** using Go channels
- **HTTP API** for enqueuing events
- **Webhook delivery**: pushes events to a configured consumer URL
- **Graceful shutdown**: ensures all queued messages are delivered before exit
- **Configurable** via environment variables or CLI flags
- **Docker-ready** for easy deployment
- **Cloudevents-compatible** message format

---

## Quick Start

### Run with Docker

```bash
docker run \
	-e CAPACITY=1000 \
	-e CONSUMER_URL=http://localhost:4000 \
	-p 3000:3000 \
	--name queue github.com/nicograef/queue
```

### Run Locally

```sh
CAPACITY=1000 CONSUMER_URL=http://localhost:4000 go run .
```

---

## Configuration

You can configure Qugo using environment variables or CLI flags:

| Variable        | Default                  | Description                       |
|----------------|-------------------------|-----------------------------------|
| `PORT`         | `3000`                   | Port for HTTP server              |
| `CAPACITY`     | `1000`                   | Max number of queued messages     |
| `CONSUMER_URL` | `http://localhost:4000`  | Webhook URL for event delivery    |

---

## API Documentation

### Enqueue event message

**POST /**

**Content-Type:** `application/json`

**Payload Example:**

```json
{
  "id": "b8e7c2e2-1f4a-4c2e-9c3a-8f7d2b6e4a1f",
  "type": "com.example.event:v1",
  "time": "2025-09-14T12:34:56Z",
  "source": "https://example.com",
  "subject": "/users/12345",
  "data": { "payload": "this is some data" }
}
```

**Response:**

```json
{ "ok": true, "queueSize": 1 }
```

---

## Development

### Build & Run

```sh
go build -o queue .
./queue
```

### Build Docker Image

```sh
docker build -t github.com/nicograef/queue .
```

### Run Docker Container

```sh
docker run -p 3000:3000 github.com/nicograef/queue
```

### Run Tests

```sh
go test ./...
```

---

## Example: Enqueue a Message

```sh
curl -X POST http://localhost:3000 \
    -H "Content-Type: application/json" \
    -d '{
        "id": "b8e7c2e2-1f4a-4c2e-9c3a-8f7d2b6e4a1f",
        "type": "com.example.event:v1",
        "time": "2024-06-13T00:00:00Z",
        "source": "https://example.com",
        "subject": "/users/12345",
        "data": { "payload": "this is some data" }
    }'
```

---

## Related Projects

- [Cloudevents spec](https://github.com/cloudevents/spec/blob/v1.0.2/cloudevents/spec.md)
- [EventsourcingDB by @thenativeweb](https://www.thenativeweb.io/products/eventsourcingdb)

---

## License

MIT
