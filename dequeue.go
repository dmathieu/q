package q

// EventHandler is a generic interface to handle events
// When an event is dequeued, it's `Handle` method will be called
type EventHandler interface {
	Handle([]byte) error
}

// Dequeue fetches an enqueued record and processes it
func (q *Queue) Dequeue(handler EventHandler) error {
	return nil
}
