package core

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

// SendToWebhook posts the message to the consumer webhook and returns the response body or error
func SendToWebhook(url string, msg Message) (string, error) {
	importBytes, err := json.Marshal(msg)
	if err != nil {
		return "", err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(importBytes))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
