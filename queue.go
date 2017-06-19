package q

import "errors"

// A Datastore allows communicating with a data storage for storing and retrieving record
type Datastore interface {
	Store([]byte) error
	Retrieve() ([]byte, error)
	Length() (int, error)
}

// A Queue allows enqueuing and listening to events
type Queue struct {
	store Datastore
}

// DataStore is used as an argument to `queue.New` to set a custom datastore
func DataStore(s Datastore) func(q *Queue) error {
	return func(q *Queue) error {
		q.store = s
		return nil
	}
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

	return handler(r)
}
