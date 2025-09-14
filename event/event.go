package event

import (
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Event represents an event (message) following the CNCF Cloudevents specification.
type Event struct {
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

func New(eventType, source, subject string, data any) (*Event, error) {
	event := Event{
		ID:      uuid.New(),
		Type:    eventType,
		Time:    time.Now().UTC(),
		Source:  source,
		Subject: subject,
		Data:    data,
	}

	if err := event.Validate(); err != nil {
		return nil, err
	}

	return &event, nil
}

func FromJSON(s string) (*Event, error) {
	var event Event

	if err := json.Unmarshal([]byte(s), &event); err != nil {
		return nil, err
	}

	if err := event.Validate(); err != nil {
		return nil, err
	}

	return &event, nil
}

// Validate checks the Event fields for validity according to the CNCF Cloudevents specification.
func (e *Event) Validate() error {
	if e.ID == uuid.Nil {
		return errors.New("event ID cannot be nil")
	}

	if len(strings.TrimSpace(e.Type)) < 5 {
		return errors.New("event type must be at least 5 characters long")
	}

	if e.Time.IsZero() {
		return errors.New("event time cannot be zero")
	}

	if len(strings.TrimSpace(e.Source)) < 5 {
		return errors.New("event source must be at least 5 characters long")
	}
	if !strings.HasPrefix(e.Source, "http://") && !strings.HasPrefix(e.Source, "https://") {
		return errors.New("event source must be a valid URI starting with http:// or https://")
	}

	if len(strings.TrimSpace(e.Subject)) < 5 {
		return errors.New("event subject must be at least 5 characters long")
	}

	if e.Data == nil {
		return errors.New("event data cannot be nil")
	}

	return nil
}
