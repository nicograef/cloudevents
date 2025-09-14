package database

import (
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/nicograef/cloudevents/event"
)

func TestPersistToJsonFile(t *testing.T) {
	defer os.Remove("database.json") // Clean up after test

	db := New()
	if db == nil {
		t.Fatal("Failed to create database")
	}

	_, err := db.AddEvent(event.EventCandidate{
		Type:    "user.new",
		Source:  "https://example.com",
		Subject: "/users/1",
		Data:    user{"ID": "1", "Name": "John Doe"},
	})
	if err != nil {
		t.Fatal("Failed to add event:", err)
	}

	if err := db.PersistToJsonFile(); err != nil {
		t.Fatal("Failed to persist to JSON file:", err)
	}
}

func TestLoadDatabaseFromJsonFile(t *testing.T) {
	// Write a sample JSON file with the correct format (array of events)
	eventJSON := `[{"id":"f8ceae97-5a98-473d-a075-c1f0a530da2c","type":"user.new","time":"2025-09-01T17:09:53.68409515Z","source":"https://example.com","subject":"/users/1","data":{"ID":"1","Name":"John Doe"}}]`
	err := os.WriteFile("database.json", []byte(eventJSON), 0644)
	if err != nil {
		t.Fatal("Failed to write database.json:", err)
	}
	defer os.Remove("database.json")

	db, err := LoadFromJSONFile()
	if err != nil {
		t.Fatal("Failed to load database from JSON file:", err)
	}

	if db == nil {
		t.Fatal("Loaded database is nil")
	}

	if len(db.Events) != 1 {
		t.Fatal("Failed to load events")
	}

	id, err := uuid.Parse("f8ceae97-5a98-473d-a075-c1f0a530da2c")
	if err != nil {
		t.Fatal("Failed to parse UUID:", err)
	}
	event := db.GetEvent(id)
	if event == nil {
		t.Fatal("Failed to get event by ID")
	}

	// Convert the event data to user map for comparison
	eventData, ok := event.Data.(map[string]interface{})
	if !ok {
		t.Fatal("Event data is not a map")
	}

	expectedUser := user{"ID": "1", "Name": "John Doe"}
	if !equalUser(eventData, expectedUser) {
		t.Fatalf("Event data mismatch. Got: %v, Expected: %v", eventData, expectedUser)
	}
}
