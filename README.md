# CloudEvents Toolkit

A comprehensive toolkit for building event-driven applications in Go, following the [CNCF CloudEvents specification](https://cloudevents.io). This monorepo contains three complementary modules that work together to provide a complete event sourcing and messaging solution.

---

## 📦 Modules

### [📚 Event](./event/)

A lightweight Go library providing CloudEvents-compatible event types and validation.

- ✅ CloudEvents-compliant `Event` struct
- 🔍 Built-in validation with clear error messages
- 🔧 Safe constructors and JSON parsing helpers
- 📋 Minimal dependencies (only UUID generation)

```go
event, err := event.New(event.EventCandidate{
    Type:    "com.example.user.created:v1",
    Source:  "https://api.example.com",
    Subject: "/users/123",
    Data:    map[string]any{"name": "John Doe"},
})
```

### [🗄️ Database](./database/)

An in-memory event-sourcing database with persistence and indexing capabilities.

- 💾 Fast in-memory event storage with JSON persistence
- 🔎 Indexed queries by event type and subject
- 🌐 HTTP API for external integrations
- 🐳 Docker-ready with volume mounting support
- 🛡️ Graceful shutdown with automatic data persistence

```bash
curl -X POST http://localhost:5000/add \
    -H "Content-Type: application/json" \
    -d '{"type": "user.created", "source": "api", "subject": "/users/123", "data": {...}}'
```

### [📨 Queue](./queue/)

An asynchronous message queue for reliable event delivery to webhooks.

- ⚡ High-performance Go channel-based queueing
- 🔄 Reliable webhook delivery with retry logic
- 📊 Configurable capacity and delivery settings
- 🌐 HTTP API for event submission
- 🐳 Production-ready Docker deployment

```bash
curl -X POST http://localhost:3000 \
    -H "Content-Type: application/json" \
    -d '{"type": "order.created", "source": "shop", "subject": "/orders/456", "data": {...}}'
```

---

## 🚀 Quick Start

### Option 1: Use Individual Modules

Each module can be used independently in your Go projects:

```bash
# Install the event library
go get github.com/nicograef/cloudevents/event

# Install the database module
go get github.com/nicograef/cloudevents/database

# Install the queue module
go get github.com/nicograef/cloudevents/queue
```

### Option 2: Run with Docker Compose

Create a complete event-driven system with all three components:

```yaml
# docker-compose.yml
version: "3.8"
services:
  database:
    image: github.com/nicograef/cloudevents/database
    ports:
      - "5000:5000"
    environment:
      - DATA_DIR=/data
    volumes:
      - ./data:/data

  queue:
    image: github.com/nicograef/cloudevents/queue
    ports:
      - "3000:3000"
    environment:
      - CAPACITY=1000
      - CONSUMER_URL=http://your-webhook-endpoint
```

```bash
docker-compose up -d
```

### Option 3: Development Setup

Clone and run locally for development:

```bash
git clone https://github.com/nicograef/cloudevents.git
cd cloudevents

# Run the database
cd database && go run .

# Run the queue (in another terminal)
cd queue && go run .
```

---

## 🏗️ Architecture

The modules work together to provide a complete event-driven architecture:

```
┌─────────────┐    HTTP POST     ┌─────────────┐    Webhooks    ┌─────────────┐
│   Client    │ ──────────────► │    Queue    │ ─────────────► │  Consumer   │
│ Application │                 │             │                │  Services   │
└─────────────┘                 └─────────────┘                └─────────────┘
       │                               │
       │ HTTP POST                     │ Optional: Store events
       ▼                               ▼
┌─────────────┐                 ┌─────────────┐
│  Database   │                 │  Database   │
│   (Events)  │                 │   (Audit)   │
└─────────────┘                 └─────────────┘
```

**Event Flow:**

1. **Clients** submit events to the **Queue** via HTTP
2. **Queue** delivers events to configured webhook endpoints
3. **Database** stores events for querying and audit trails
4. **Event library** ensures consistent CloudEvents format across all components

---

## 🔧 Configuration

Each module supports environment-based configuration:

| Module   | Variable       | Default                 | Description                |
| -------- | -------------- | ----------------------- | -------------------------- |
| Database | `PORT`         | `5000`                  | HTTP server port           |
| Database | `DATA_DIR`     | `.`                     | Data persistence directory |
| Queue    | `PORT`         | `3000`                  | HTTP server port           |
| Queue    | `CAPACITY`     | `1000`                  | Max queued messages        |
| Queue    | `CONSUMER_URL` | `http://localhost:4000` | Webhook delivery endpoint  |

---

## 📖 Examples

### Basic Event Creation and Storage

```go
package main

import (
    "fmt"
    "github.com/nicograef/cloudevents/database/database"
    "github.com/nicograef/cloudevents/event"
)

func main() {
    // Create event
    candidate := event.EventCandidate{
        Type:    "com.example.user.signup:v1",
        Source:  "https://myapp.com",
        Subject: "/users/123",
        Data:    map[string]any{"email": "user@example.com"},
    }

    // Store in database
    db := database.New()
    storedEvent, err := db.AddEvent(candidate)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Event stored with ID: %s\n", storedEvent.ID)

    // Query events
    userEvents := db.GetEventsBySubject("/users/123")
    fmt.Printf("Found %d events for user\n", len(userEvents))
}
```

### HTTP Event Submission

```bash
# Submit to queue for delivery
curl -X POST http://localhost:3000 \
    -H "Content-Type: application/json" \
    -d '{
        "type": "com.shop.order.created:v1",
        "source": "https://shop.example.com",
        "subject": "/orders/12345",
        "data": {"amount": 99.99, "currency": "USD"}
    }'

# Store in database for querying
curl -X POST http://localhost:5000/add \
    -H "Content-Type: application/json" \
    -d '{
        "type": "com.shop.order.created:v1",
        "source": "https://shop.example.com",
        "subject": "/orders/12345",
        "data": {"amount": 99.99, "currency": "USD"}
    }'
```

---

## 🧪 Testing

Run tests across all modules:

```bash
# Test all modules
find . -name "go.mod" -execdir go test ./... \;

# Test individual modules
cd event && go test ./...
cd database && go test ./...
cd queue && go test ./...
```

---

## 🚢 Production Deployment

### Kubernetes

Example deployment manifests:

```yaml
# database-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: cloudevents-database
spec:
  replicas: 1
  selector:
    matchLabels:
      app: cloudevents-database
  template:
    metadata:
      labels:
        app: cloudevents-database
    spec:
      containers:
        - name: database
          image: github.com/nicograef/cloudevents/database
          ports:
            - containerPort: 5000
          env:
            - name: DATA_DIR
              value: /data
          volumeMounts:
            - name: data-volume
              mountPath: /data
      volumes:
        - name: data-volume
          persistentVolumeClaim:
            claimName: database-pvc
```

### Docker Swarm

```yaml
# docker-stack.yml
version: "3.8"
services:
  database:
    image: github.com/nicograef/cloudevents/database
    deploy:
      replicas: 1
    environment:
      DATA_DIR: /data
    volumes:
      - database-data:/data
    ports:
      - "5000:5000"

  queue:
    image: github.com/nicograef/cloudevents/queue
    deploy:
      replicas: 3
    environment:
      CAPACITY: 5000
      CONSUMER_URL: http://your-webhook-service
    ports:
      - "3000:3000"

volumes:
  database-data:
```

---

## 🤝 Contributing

We welcome contributions! Here's how to get started:

1. **Fork** the repository
2. **Clone** your fork: `git clone https://github.com/yourusername/cloudevents.git`
3. **Create** a feature branch: `git checkout -b feature/amazing-feature`
4. **Make** your changes and add tests
5. **Test** everything: `find . -name "go.mod" -execdir go test ./... \;`
6. **Commit** your changes: `git commit -m 'Add amazing feature'`
7. **Push** to your branch: `git push origin feature/amazing-feature`
8. **Open** a Pull Request

### Development Guidelines

- Follow Go best practices and `gofmt` formatting
- Add tests for new functionality
- Update documentation and examples
- Ensure backwards compatibility when possible
- Keep modules loosely coupled

---

## 📋 Roadmap

- [ ] **Clustering support** for horizontal scaling
- [ ] **Event replay** capabilities in database
- [ ] **Dead letter queues** in queue module
- [ ] **Metrics and observability** endpoints
- [ ] **gRPC APIs** alongside HTTP
- [ ] **Stream processing** capabilities
- [ ] **Event schema registry** integration

---

## 📚 Related Projects

- [CloudEvents Specification](https://github.com/cloudevents/spec) - Official CNCF specification
- [CloudEvents SDK](https://github.com/cloudevents/sdk-go) - Official Go SDK
- [EventStore](https://www.eventstore.com/) - Production event sourcing database
- [Apache Kafka](https://kafka.apache.org/) - Distributed event streaming platform
- [NATS](https://nats.io/) - Cloud native messaging system

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙋‍♂️ Support

- 📖 **Documentation**: Check individual module READMEs for detailed usage
- 🐛 **Bug Reports**: [Open an issue](https://github.com/nicograef/cloudevents/issues)
- 💡 **Feature Requests**: [Start a discussion](https://github.com/nicograef/cloudevents/discussions)
- 📧 **Contact**: [nico@example.com](mailto:nico@example.com)

---

<div align="center">

**Built with ❤️ for the event-driven future**

[⭐ Star this repo](https://github.com/nicograef/cloudevents) | [🍴 Fork it](https://github.com/nicograef/cloudevents/fork) | [📖 Docs](https://github.com/nicograef/cloudevents#modules)

</div>
