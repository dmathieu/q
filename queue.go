package q

import (
	"errors"

	"github.com/dmathieu/q/stores"
)

// A Queue allows enqueuing and listening to events
type Queue struct {
	store stores.Datastore
}

// New initializes a new queue, with options
func New(options ...func(*Queue) error) (*Queue, error) {
	q := &Queue{}

	for _, option := range options {
		if err := option(q); err != nil {
			return nil, err
		}
	}

	if q.store == nil {
		return nil, errors.New("no data store specified")
	}

	return q, nil
}

// Enqueue enqueues a new entry to be processed
func (q *Queue) Enqueue(c []byte) error {
	return q.store.Store(c)
}

// Handle fetches an enqueued record and processes it
func (q *Queue) Handle(handler func([]byte) error) error {
	r, err := q.store.Retrieve()
	if err != nil {
		return nil
	}

	if r == nil {
		return nil
	}

	// TODO if we get an error in finish, we lose the potential handler error
	err = handler(r)
	err2 := q.store.Finish(r)
	if err2 != nil {
		return err2
	}
	return err
}
