package api

import (
	"log"
	"net/http"

	"github.com/nicograef/cloudevents/event"
	"github.com/nicograef/cloudevents/queue/queue"
)

// EnqueueResponseSuccess represents a successful response from the enqueue API endpoint.
type EnqueueResponseSuccess struct {
	Ok        bool `json:"ok"`
	QueueSize int  `json:"queueSize"`
}

// EnqueueResponseError represents a failed response from the enqueue API endpoint.
type EnqueueResponseError struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error"`
}

// NewEnqueueHandler returns an HTTP handler for enqueuing messages into the queue.
// It expects a POST request with a JSON body containing a 'message' field.
func NewEnqueueHandler(appQueue queue.Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		message := event.Event{}
		if !readJSONRequest(w, r, &message) {
			return
		}

		if err := message.Validate(); err != nil {
			log.Printf("Invalid event: %v", err)
			sendJSONResponse(w, EnqueueResponseError{
				Ok:    false,
				Error: err.Error(),
			})
			return
		}

		appQueue.Queue <- queue.QueueMessage{Message: message, Attempts: 0}

		sendJSONResponse(w, EnqueueResponseSuccess{
			Ok:        true,
			QueueSize: len(appQueue.Queue),
		})
	}
}
