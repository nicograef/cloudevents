package database

import (
	"testing"

	"github.com/google/uuid"
	"github.com/nicograef/cloudevents/event"
)

type user map[string]any

func equalUser(a, b user) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}

func TestDatabaseCreation(t *testing.T) {
	db := New()

	if db == nil {
		t.Fatal("Failed to create database")
	}
	if db.Events == nil {
		t.Fatal("Failed to create events")
	}
	if len(db.Events) != 0 {
		t.Fatal("Failed to create events")
	}
}

func TestAddAndGetEvents(t *testing.T) {
	db := New()
	event1, err := db.AddEvent(event.Candidate{Type: "user.new", Source: "https://example.com", Subject: "/users/1", Data: user{"ID": "1", "Name": "John Doe"}})
	if err != nil {
		t.Fatalf("AddEvent failed: %v", err)
	}
	event2, err := db.AddEvent(event.Candidate{Type: "user.update", Source: "https://example.com", Subject: "/users/1", Data: user{"ID": "1", "Email": "john.doe@example.com"}})
	if err != nil {
		t.Fatalf("AddEvent failed: %v", err)
	}

	if len(db.Events) != 2 {
		t.Fatal("Failed to add events")
	}
	if len(db.TypeIndex) != 2 {
		t.Fatal("Failed to create type index")
	}
	if len(db.SubjectIndex) != 1 {
		t.Fatal("Failed to create entity index")
	}

	events := db.GetEvents()
	if len(events) != 2 {
		t.Fatal("Failed to get events")
	}

	if event1.ID != events[0].ID {
		t.Fatal("Event ID is not the same as the one created")
	}
	if event2.ID != events[1].ID {
		t.Fatal("Event ID is not the same as the one created")
	}

	if !equalUser(events[0].Data.(user), user{"ID": "1", "Name": "John Doe"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}
	if !equalUser(events[1].Data.(user), user{"ID": "1", "Email": "john.doe@example.com"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}

	if events[0].Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}
	if events[1].Type != "user.update" {
		t.Fatal("Event type is not the same as the one created")
	}
}

func TestGetEventByID(t *testing.T) {
	db := New()
	event1, err := db.AddEvent(event.Candidate{Type: "user.new", Source: "https://example.com", Subject: "/users/1", Data: user{"ID": "1", "Name": "John Doe"}})
	if err != nil {
		t.Fatalf("AddEvent failed: %v", err)
	}

	if nonExistingEvent := db.GetEvent(uuid.New()); nonExistingEvent != nil {
		t.Fatal("Expected no event to be found")
	}

	event := db.GetEvent(event1.ID)
	if event == nil {
		t.Fatal("Failed to get event by ID")
	}

	if !equalUser(event.Data.(user), user{"ID": "1", "Name": "John Doe"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}
	if event.Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}
	if err := uuid.Validate(event.ID.String()); err != nil {
		t.Fatal("Event ID is not valid")
	}
}

func TestGetEventsByType(t *testing.T) {
	db := New()
	db.AddEvent(event.Candidate{Type: "user.new", Source: "https://example.com", Subject: "/users/1", Data: user{"ID": "1", "Name": "John Doe"}})
	db.AddEvent(event.Candidate{Type: "user.update", Source: "https://example.com", Subject: "/users/1", Data: user{"ID": "1", "Email": "john.doe@example.com"}})
	db.AddEvent(event.Candidate{Type: "user.new", Source: "https://example.com", Subject: "/users/2", Data: user{"ID": "2", "Name": "Max Mustermann"}})

	if nonExistingEvents := db.GetEventsByType("non-existing-type"); len(nonExistingEvents) != 0 {
		t.Fatal("Expected no event to be found")
	}

	events := db.GetEventsByType("user.new")
	if len(events) != 2 {
		t.Fatal("Failed to get events by type")
	}

	if events[0].Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}
	if events[1].Type != "user.new" {
		t.Fatal("Event type is not the same as the one created")
	}

	if !equalUser(events[0].Data.(user), user{"ID": "1", "Name": "John Doe"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}
	if !equalUser(events[1].Data.(user), user{"ID": "2", "Name": "Max Mustermann"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}
}

func TestGetEventsBySubject(t *testing.T) {
	db := New()
	db.AddEvent(event.Candidate{Type: "user.new", Source: "https://example.com", Subject: "/users/1", Data: user{"ID": "1", "Name": "John Doe"}})
	db.AddEvent(event.Candidate{Type: "user.update", Source: "https://example.com", Subject: "/users/1", Data: user{"ID": "1", "Email": "john.doe@example.com"}})
	db.AddEvent(event.Candidate{Type: "user.new", Source: "https://example.com", Subject: "/users/2", Data: user{"ID": "2", "Name": "Max Mustermann"}})

	if nonExistingEvents := db.GetEventsBySubject("/non-existing"); len(nonExistingEvents) != 0 {
		t.Fatal("Expected no event to be found")
	}

	events := db.GetEventsBySubject("/users/1")
	if len(events) != 2 {
		t.Fatal("Failed to get events by subject")
	}

	if events[0].Subject != "/users/1" {
		t.Fatal("Event subject is not the same as the one created")
	}
	if events[1].Subject != "/users/1" {
		t.Fatal("Event subject is not the same as the one created")
	}

	if !equalUser(events[0].Data.(user), user{"ID": "1", "Name": "John Doe"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}
	if !equalUser(events[1].Data.(user), user{"ID": "1", "Email": "john.doe@example.com"}) {
		t.Fatal("Event retrieved is not the same as the one created")
	}
}
