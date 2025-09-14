package queue

import (
	"log"

	"github.com/nicograef/cloudevents/event"
)

type QueueMessage struct {
	Message  event.Event
	Attempts int
}

type Queue struct {
	Queue       chan QueueMessage
	FailedQueue []QueueMessage
}

func NewQueue(capacity int) Queue {
	return Queue{Queue: make(chan QueueMessage, capacity), FailedQueue: []QueueMessage{}}
}

// StartConsumer starts a goroutine that reads from the queue and calls the webhook for each message.
// It takes the queue, consumerURL, and a WaitGroup pointer.
// SendFunc defines the signature for sending a message to a webhook.
type SendFunc func(url string, msg event.Event) (string, error)

func (q *Queue) HandleQueueItem(item QueueMessage, consumerURL string, sendFunc SendFunc) {
	resp, err := sendFunc(consumerURL, item.Message)

	if err != nil {
		log.Printf("Error sending to webhook: %v", err)

		item.Attempts++

		if item.Attempts < 3 {
			log.Printf("Re-enqueuing message, attempt %d", item.Attempts)
			q.Queue <- item
		} else {
			log.Printf("Max attempts reached for message: %+v", item.Message)
			q.FailedQueue = append(q.FailedQueue, item)
		}
	} else {
		log.Printf("Webhook response: %s", resp)
	}
}
