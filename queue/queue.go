package queue

import (
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"
)

var ErrInvalidState = errors.New("invalid state")
var ErrQueueEmpty = errors.New("queue is empty")

type Queue struct {
	Messages []Message `json:"messages"`
}

func New() *Queue {
	return &Queue{
		Messages: make([]Message, 0),
	}
}

// Peek returns the first message without removing it from the queue.
// If the queue is empty, it returns nil.
func (q *Queue) Peek() *Message {
	if q.Size() == 0 {
		return nil
	}

	return &q.Messages[0]
}

// Add a new message to the end of the queue.
func (q *Queue) Enqueue(message Message) {
	q.Messages = append(q.Messages, message)
}

func (q *Queue) Dequeue(messageId uuid.UUID) error {
	if q.Size() == 0 {
		return ErrQueueEmpty
	}

	message := q.Messages[0]
	if message.ID != messageId {
		return ErrInvalidState
	}

	// Remove the first message from the queue
	q.Messages = q.Messages[1:]

	return nil
}

func (q *Queue) Size() int {
	return len(q.Messages)
}

func (q *Queue) PersistToJsonFile() error {
	data, err := json.MarshalIndent(q, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("queue.json", data, 0644)
}

func LoadQueueFromJsonFile() (*Queue, error) {
	data, err := os.ReadFile("queue.json")
	if err != nil {
		return nil, err
	}

	var q Queue
	if err := json.Unmarshal(data, &q); err != nil {
		return nil, err
	}

	return &q, nil
}
