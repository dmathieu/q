package q

import (
	"errors"
	"fmt"
	"time"

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

// HouseKeeping handles regular housekeeping tasks the datastore needs to perform
func (q *Queue) HouseKeeping() error {
	return q.store.HouseKeeping()
}

// Run starts a local worker
func (q *Queue) Run(handler func([]byte) error, mc int) error {
	c := make(chan struct{}, mc)
	errCh := make(chan error)

	go func() {
		defer func() {
			if x := recover(); x != nil {
				err, ok := x.(error)
				if !ok {
					err = fmt.Errorf("%q", err)
				}
				errCh <- err
			}
		}()

		for {
			select {
			case <-time.After(time.Minute):
				err := q.HouseKeeping()
				if err != nil {
					errCh <- err
				}
			}
		}
	}()

	for {
		select {
		case err := <-errCh:
			return err
		case c <- struct{}{}:
			go func() {
				defer func() {
					if x := recover(); x != nil {
						err, ok := x.(error)
						if !ok {
							err = fmt.Errorf("%q", err)
						}
						errCh <- err
					}

					<-c
				}()

				err := q.Handle(handler)
				if err != nil {
					errCh <- err
				}
			}()
		}
	}
}
