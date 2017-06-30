package q

import (
	"errors"

	"github.com/dmathieu/q/stores"
)

// A Queue allows enqueuing and listening to events
type Queue struct {
	store          stores.Datastore
	failureHandler func([]byte) error
}

// FailureHandler is an option for new queues, to set a method executed when a record failes executing
func FailureHandler(f func([]byte) error) func(q *Queue) error {
	return func(q *Queue) error {
		q.failureHandler = f
		return nil
	}
}

// New initializes a new queue, with options
func New(options ...func(*Queue) error) (*Queue, error) {
	q := &Queue{
		failureHandler: func([]byte) error { return nil },
	}

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

	err = handler(r)
	if err != nil {
		// TODO if we get an error in finish, we lose the potential failure error
		err = q.failureHandler(r)
	}
	err2 := q.store.Finish(r)
	if err2 != nil {
		return err2
	}
	return err
}

// HouseKeeping handles regular housekeeping tasks the datastore needs to perform
func (q *Queue) HouseKeeping() error {
	return q.store.HouseKeeping()
}
