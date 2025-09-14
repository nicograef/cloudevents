package api

import (
	"net/http"

	"github.com/nicograef/qugo/core"
)

// EnqueueResponse represents the response from the enqueue API endpoint.
type EnqueueResponse struct {
	Ok        bool `json:"ok"`
	QueueSize int  `json:"queueSize"`
}

// NewEnqueueHandler returns an HTTP handler for enqueuing messages into the queue.
// It expects a POST request with a JSON body containing a 'message' field.
func NewEnqueueHandler(appQueue chan core.Message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !validateMethod(w, r, http.MethodPost) {
			return
		}

		message := core.Message{}
		if !readJSONRequest(w, r, &message) {
			return
		}

		appQueue <- message

		sendJSONResponse(w, EnqueueResponse{
			Ok:        true,
			QueueSize: len(appQueue),
		})
	}
}
