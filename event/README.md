# CloudEvents Event (Go)

A tiny Go library that provides a simple Event type and helpers aligned with the CNCF CloudEvents model. It includes:

- An `Event` struct with common CloudEvents-style attributes
- A safe constructor `New(...)` that validates inputs
- A `Validate()` method you can call on any event
- A `FromJSON(...)` helper to parse and validate JSON payloads

Module path: `github.com/nicograef/cloudevents/event`

## Install

```bash
# From another module
# Adds this package as a dependency
# (If you are developing locally, see the replace tip below.)
go get github.com/nicograef/cloudevents/event
```

If you are developing locally in a monorepo alongside this module, you can point your consumer module to the local folder using a replace directive in its `go.mod`:

```go
replace github.com/nicograef/cloudevents/event => ../event
```

## Quick start

```go
package main

import (
    "encoding/json"
    "fmt"

    event "github.com/nicograef/cloudevents/event"
)

func main() {
    // Create a new event (validated)
    e, err := event.New(
        "com.example.order.created:v1",
        "https://shop.example.com",
        "/orders/42",
        map[string]any{"amount": 19.99, "currency": "USD"},
    )
    if err != nil {
        panic(err)
    }

    // Marshal to JSON
    b, _ := json.Marshal(e)
    fmt.Println(string(b))

    // Parse back from JSON (validated)
    parsed, err := event.FromJSON(string(b))
    if err != nil {
        panic(err)
    }
    fmt.Println("parsed type:", parsed.Type)
}
```

## API

- type `Event` struct
  - `ID uuid.UUID` — unique event ID (auto-assigned by `New`)
  - `Type string` — e.g. `com.example.something:v1`
  - `Time time.Time` — UTC timestamp (auto-set by `New`)
  - `Source string` — URI identifying the producer, e.g. `https://service.example.com`
  - `Subject string` — entity or resource within the source, e.g. `/users/123`
  - `Data any` — event payload (any JSON-marshalable value)

- func `New(eventType, source, subject string, data any) (*Event, error)`
  - Creates an `Event` with generated `ID` and current UTC `Time`, then validates it.

- func `FromJSON(s string) (*Event, error)`
  - Unmarshals JSON into an `Event` and validates it.

- method `(e *Event) Validate() error`
  - Validates the event fields (see rules below).

## Validation rules

`Validate()` enforces the following:

- `ID` must be non-nil
- `Type` must be at least 5 characters
- `Time` cannot be zero
- `Source` must be at least 5 characters and start with `http://` or `https://`
- `Subject` must be at least 5 characters
- `Data` cannot be nil

These checks are run in `New(...)` and `FromJSON(...)`, and you can call `Validate()` manually after any mutation.

## JSON shape (example)

```json
{
  "id": "550e8400-e29b-41d4-a716-446655440000",
  "type": "com.example.event:v1",
  "time": "2023-01-01T12:00:00Z",
  "source": "https://example.com",
  "subject": "/users/123",
  "data": {"key": "value"}
}
```

## Testing

From this module's directory (`event/`):

```bash
go test ./...
```

## Notes

- This library aims to be lightweight and practical while keeping familiar CloudEvents semantics. It does not attempt to implement the entire CloudEvents spec; instead, it provides a minimal, validated event shape that works well for many services.
- If you need stricter conformance or protocol bindings, consider the official CloudEvents SDKs.
