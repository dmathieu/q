package q

// Dequeue fetches an enqueued record and processes it
func (q *Queue) Dequeue(handler func([]byte) error) error {
	r, err := q.store.Retrieve()
	if err != nil {
		return nil
	}

	if r == nil {
		return nil
	}

	return handler(r)
}
