package database

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"

	"github.com/google/uuid"
	"github.com/nicograef/cloudevents/event"
)

// LoadFromJSONFile loads the database state from a JSON file on disk.
// If the file does not exist or cannot be read, an error is returned.
// The indexes are rebuilt after loading the events.
func LoadFromJSONFile(dataDir string) (*Database, error) {
	filePath := filepath.Join(dataDir, "database.json")
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Println("Error closing file:", closeErr)
		}
	}()

	// Decode the JSON data into a slice of events
	var events []event.Event
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&events); err != nil {
		return nil, err
	}

	// Rebuild the events map from the loaded events slice
	eventsMap := make(map[uuid.UUID]event.Event)
	for _, e := range events {
		eventsMap[e.ID] = e
	}

	db := Database{Events: eventsMap}
	db.RebuildIndexes()

	return &db, nil
}

// PersistToJsonFile saves the current state of the database to the disk.
// The events are stored as an array in a JSON format for easy parsing.
// The indexes are not persisted to save space and can be rebuilt on load.
func (db *Database) PersistToJsonFile(dataDir string) error {
	filePath := filepath.Join(dataDir, "database.json")
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Println("Error closing file:", closeErr)
		}
	}()

	// convert map to slice for easier JSON encoding
	events := make([]event.Event, 0, len(db.Events))
	for _, e := range db.Events {
		events = append(events, e)
	}

	// Encode the events slice to JSON and write to file
	encoder := json.NewEncoder(file)
	if err := encoder.Encode(events); err != nil {
		return err
	}

	return nil
}
