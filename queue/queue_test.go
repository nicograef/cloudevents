package queue

import (
	"testing"
)

func TestHandleQueueItem_Success(t *testing.T) {
	q := NewQueue(1)
	called := false
	sendFunc := func(url string, msg Message) (string, error) {
		called = true
		return "ok", nil
	}
	item := QueueMessage{Message: Message{Type: "success"}, Attempts: 0}
	q.HandleQueueItem(item, "http://test", sendFunc)
	if !called {
		t.Errorf("sendFunc was not called")
	}
	select {
	case <-q.Queue:
		t.Errorf("queue should be empty after success")
	default:
	}
	if len(q.FailedQueue) != 0 {
		t.Errorf("FailedQueue should be empty after success")
	}
}

func TestHandleQueueItem_RetryAndMaxAttempts(t *testing.T) {
	q := NewQueue(1)
	attempts := 0
	sendFunc := func(url string, msg Message) (string, error) {
		attempts++
		return "", assertError("fail")
	}
	item := QueueMessage{Message: Message{Type: "fail"}, Attempts: 0}
	q.HandleQueueItem(item, "http://test", sendFunc)
	// Should be re-enqueued with Attempts=1
	var retried QueueMessage
	select {
	case retried = <-q.Queue:
		if retried.Attempts != 1 {
			t.Errorf("expected Attempts=1, got %d", retried.Attempts)
		}
		// Try again, should be re-enqueued with Attempts=2
		q.HandleQueueItem(retried, "http://test", sendFunc)
		retried2 := <-q.Queue
		if retried2.Attempts != 2 {
			t.Errorf("expected Attempts=2, got %d", retried2.Attempts)
		}
		// Try again, should NOT be re-enqueued (Attempts=3), should go to FailedQueue
		q.HandleQueueItem(retried2, "http://test", sendFunc)
		select {
		case <-q.Queue:
			t.Errorf("should not re-enqueue after max attempts")
		default:
		}
		if len(q.FailedQueue) != 1 {
			t.Errorf("expected 1 failed message, got %d", len(q.FailedQueue))
		}
	default:
		t.Errorf("expected message to be re-enqueued")
	}
}

// assertError is a helper to create a test error
func assertError(msg string) error {
	return &testError{msg}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
