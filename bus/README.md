# bus

**bus** is a message bus for event-driven applications, written in Go. It follows the [CNCF Cloudevents specification](https://cloudevents.io) and is designed to receive event messages via HTTP and publish them to multple subscribers (event message receivers defined as webhooks).

---

## Quick Start

### Run with Docker

```bash
docker run \
	-e SUBSCRIBER_URLS=http://localhost:4000 \
	-p 3000:3000 \
	--name bus github.com/nicograef/bus
```

### Run Locally

```sh
SUBSCRIBER_URLS=http://localhost:4000 go run .
```

---

## Configuration

You can configure **bus** using environment variables or CLI flags:

| Variable        | Default                                       | Description                     |
| --------------- | --------------------------------------------- | ------------------------------- |
| `PORT`          | `3000`                                        | Port for HTTP server            |
| `SUBSCRIBER_URLS` | `http://localhost:4000,http://localhost:5000` | Webhook URLs for event delivery |

---

## API Documentation

### Publish event message

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
{ "ok": true }
```

### Example: Publish a Message

```sh
curl -X POST http://localhost:3000/publish \
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
