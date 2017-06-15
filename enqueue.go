package q

// Enqueue enqueues a new entry to be processed
func (q *Queue) Enqueue(c []byte) error {
	return q.store.Store(c)
}
