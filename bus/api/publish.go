package api

import (
	"log"
	"net/http"

	"github.com/nicograef/cloudevents/event"
)

// PublishResponseSuccess represents a successful response from the publish API endpoint.
type PublishResponseSuccess struct {
	Ok bool `json:"ok"`
}

// PublishResponseError represents a failed response from the publish API endpoint.
type PublishResponseError struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

type PublishFunc func(e event.Event) error

// NewPublishHandler returns an HTTP handler for publishing messages to the subscribers.
// It expects a POST request with a JSON body containing a cloudevent.
func NewPublishHandler(publish PublishFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !validateMethod(w, r, http.MethodPost) {
			return
		}

		message := event.Event{}
		if !readJSONRequest(w, r, &message) {
			return
		}

		if err := message.Validate(); err != nil {
			log.Printf("Invalid event: %v", err)
			sendJSONResponse(w, PublishResponseError{
				Ok:    false,
				Error: err.Error(),
			})
			return
		}

		if err := publish(message); err != nil {
			log.Printf("Error publishing message: %v", err)
			sendJSONResponse(w, PublishResponseError{
				Ok:    false,
				Error: err.Error(),
			})
			return
		}

		sendJSONResponse(w, PublishResponseSuccess{
			Ok: true,
		})
	}
}
