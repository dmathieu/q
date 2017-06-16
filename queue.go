package q

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

	return q, nil
}
