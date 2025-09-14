package qugo

import (
	"net/http"
)

type EnqueueRequest struct {
	Message Message `json:"message"`
}
type EnqueueResponse struct {
	Ok    bool   `json:"ok"`
	Error string `json:"error,omitempty"`
}

// read the payload from the request body, create a new message with a new UUID and the current timestamp, add it to the queue, and return the ID of the new message in the response.
func NewEnqueueHandler(appQueue chan Message) http.HandlerFunc {
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
