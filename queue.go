package q

import "fmt"

// A Datastore allows communicating with a data storage for storing and retrieving record
type Datastore interface{}

// A Queue allows enqueuing and listening to events
type Queue struct {
	name  string
	store Datastore
}

// DataStore is used as an argument to `queue.New` to set a custom datastore
func DataStore(s Datastore) func(q *Queue) error {
	return func(q *Queue) error {
		q.store = s
		return nil
	}
}

// New initializes a new queue, with a name and options
func New(name string, options ...func(*Queue) error) (*Queue, error) {
	q := &Queue{name: name}

	for _, option := range options {
		if err := option(q); err != nil {
			return nil, err
		}
	}

	return q, nil
}

func (q *Queue) workingQueue() string {
	return fmt.Sprintf("%s:working", q.name)
}
