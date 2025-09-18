package bus

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/nicograef/cloudevents/bus/api"
	"github.com/nicograef/cloudevents/event"
)

type SendFunc func(url string, ev event.Event) (string, error)

// NewPublish creates a PublishFunc that sends the event to all subscribers using the provided SendFunc.
// It returns an error if sending to any subscriber fails.
func NewPublish(subs []string, send SendFunc) api.PublishFunc {
	return func(ev event.Event) error {
		var failedSubs []string

		for _, sub := range subs {
			_, err := send(sub, ev)
			if err != nil {
				log.Printf("ERROR Failed to send event %s to subscriber %s: %v", ev.ID, sub, err)
				failedSubs = append(failedSubs, sub)
			}
		}

		if len(failedSubs) > 0 {
			return fmt.Errorf("Failed to publish event %s to subscribers %v", ev.ID, failedSubs)
		}

		log.Printf("INFO Successfully published event %s to all subscribers", ev.ID)

		return nil
	}
}

// SendToWebhook posts the event to the subscriber webhook and returns the response body or error
func SendToWebhook(url string, ev event.Event) (string, error) {
	importBytes, err := json.Marshal(ev)
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
