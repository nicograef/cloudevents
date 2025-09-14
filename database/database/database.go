package database

import (
	"sort"

	"github.com/google/uuid"
	"github.com/nicograef/cloudevents/event"
)

type Database struct {
	Events       map[uuid.UUID]event.Event
	TypeIndex    map[string][]uuid.UUID
	SubjectIndex map[string][]uuid.UUID
}

func New() *Database {
	return &Database{
		Events:       make(map[uuid.UUID]event.Event),
		TypeIndex:    make(map[string][]uuid.UUID),
		SubjectIndex: make(map[string][]uuid.UUID),
	}
}

// AddEvent adds a new event to the database and updates the indexes
func (db *Database) AddEvent(candidate event.EventCandidate) (*uuid.UUID, error) {
	event, err := event.New(candidate)
	if err != nil {
		return nil, err
	}

	db.Events[event.ID] = *event
	db.TypeIndex[event.Type] = append(db.TypeIndex[event.Type], event.ID)
	db.SubjectIndex[event.Subject] = append(db.SubjectIndex[event.Subject], event.ID)

	return &event.ID, nil
}

// GetEvent retrieves an event by its ID
func (db *Database) GetEvent(id uuid.UUID) *event.Event {
	event, exists := db.Events[id]

	if !exists {
		return nil
	}

	return &event
}

// GetEvents returns all events sorted by their timestamp
func (db *Database) GetEvents() []event.Event {
	events := make([]event.Event, 0, len(db.Events))

	for _, event := range db.Events {
		events = append(events, event)
	}

	sortEventsByTime(events)

	return events
}

// GetEventsByType returns all events of a specific type sorted by their timestamp
func (db *Database) GetEventsByType(eventType string) []event.Event {
	eventIDs, exists := db.TypeIndex[eventType]

	if !exists || len(eventIDs) == 0 {
		return []event.Event{}
	}

	events := make([]event.Event, 0, len(eventIDs))
	for _, id := range eventIDs {
		if event, exists := db.Events[id]; exists {
			events = append(events, event)
		}
	}

	sortEventsByTime(events)

	return events

}

// GetEventsBySubject returns all events for a specific subject sorted by their timestamp
func (db *Database) GetEventsBySubject(subject string) []event.Event {
	eventIDs, exists := db.SubjectIndex[subject]

	if !exists || len(eventIDs) == 0 {
		return []event.Event{}
	}

	events := make([]event.Event, 0, len(eventIDs))
	for _, id := range eventIDs {
		if event, exists := db.Events[id]; exists {
			events = append(events, event)
		}
	}

	sortEventsByTime(events)

	return events
}

// RebuildIndexes reconstructs the indexes from the current events in the database
func (db *Database) RebuildIndexes() {
	db.TypeIndex = make(map[string][]uuid.UUID)
	db.SubjectIndex = make(map[string][]uuid.UUID)

	for id, event := range db.Events {
		db.TypeIndex[event.Type] = append(db.TypeIndex[event.Type], id)
		db.SubjectIndex[event.Subject] = append(db.SubjectIndex[event.Subject], id)
	}
}

// sortEventsByTime sorts events by their timestamp
func sortEventsByTime(events []event.Event) {
	sort.Slice(events, func(i, j int) bool {
		return events[i].Time.Before(events[j].Time)
	})
}
