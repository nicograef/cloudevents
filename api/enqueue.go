package api

import (
	"net/http"

	"github.com/nicograef/qugo/core"
)

// EnqueueRequest represents the expected payload for the enqueue API endpoint.
type EnqueueRequest struct {
	Message core.Message `json:"message"`
}

// EnqueueResponse represents the response from the enqueue API endpoint.
type EnqueueResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

// NewEnqueueHandler returns an HTTP handler for enqueuing messages into the queue.
// It expects a POST request with a JSON body containing a 'message' field.
func NewEnqueueHandler(appQueue chan core.Message) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !validateMethod(w, r, http.MethodPost) {
			return
		}

		requestBody := EnqueueRequest{}
		if !readJSONRequest(w, r, &requestBody) {
			return
		}

		appQueue <- requestBody.Message

		sendJSONResponse(w, EnqueueResponse{
			Ok: true,
		})
	}
}
