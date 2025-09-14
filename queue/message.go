package queue

import (
	"time"

	"github.com/google/uuid"
)

// Message represents an event message in the queue, following the CNCF Cloudevents specification.
type Message struct {
	// Identifies the event. Must be unique within the scope of the producer/source.
	ID uuid.UUID `json:"id"`
	// The type of event related to the source system and subject. E.g. com.library.book.borrowed:v1
	Type string `json:"type"`
	// The timestamp of when the event occurred.
	Time time.Time `json:"time"`
	// The source of the event. Must be a valid URI-Reference. E.g. https://library.example.com
	Source string `json:"source"`
	// The subject of the event in the context of the event producer (identified by source). E.g. the entity to which the event is primarily related. E.g. /users/12345
	Subject string `json:"subject"`
	// The event payload.
	Data any `json:"data"`
}
