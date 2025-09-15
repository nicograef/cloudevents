package api

import (
	"log"
	"net/http"

	"github.com/nicograef/cloudevents/database/database"
	"github.com/nicograef/cloudevents/event"
)

// AddEventRequest represents the expected request body for the enqueue API endpoint.
type AddEventRequest = event.Candidate

// AddEventResponseSuccess represents a successful response from the enqueue API endpoint.
type AddEventResponseSuccess struct {
	Ok    bool        `json:"ok"`
	Event event.Event `json:"event"`
}

// AddEventResponseError represents a failed response from the enqueue API endpoint.
type AddEventResponseError struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

// NewAddEventHandler creates an HTTP handler for adding events to the database.
func NewAddEventHandler(db database.Database) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !validateMethod(w, r, http.MethodPost) {
			return
		}

		candidate := event.Candidate{}
		if !readJSONRequest(w, r, &candidate) {
			return
		}

		event, err := db.AddEvent(candidate)
		if err != nil {
			log.Printf("ERROR Failed to add event to database: %v", err)
			sendJSONResponse(w, AddEventResponseError{
				Ok:    false,
				Error: err.Error(),
			})
			return
		}

		log.Printf("INFO Added event to database: %s", event.ID)

		sendJSONResponse(w, AddEventResponseSuccess{
			Ok:    true,
			Event: *event,
		})
	}

}
