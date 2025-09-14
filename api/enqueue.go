package api

import (
	"net/http"

	"github.com/nicograef/qugo/queue"
)

type EnqueueRequest struct {
	Message queue.Message `json:"message"`
}
type EnqueueResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

// read the payload from the request body, create a new message with a new UUID and the current timestamp, add it to the queue, and return the ID of the new message in the response.
func NewEnqueueHandler(q *queue.Queue) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !validateMethod(w, r, http.MethodPost) {
			return
		}

		requestBody := EnqueueRequest{}
		if !readJSONRequest(w, r, &requestBody) {
			return
		}

		q.Enqueue(requestBody.Message)

		sendJSONResponse(w, EnqueueResponse{
			Ok: true,
		})
	}
}
