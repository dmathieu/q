package queue

// A Datastore allows communicating with a data storage for storing and retrieving record
type Datastore interface {
	Store([]byte) error
	Retrieve() ([]byte, error)
	Finish([]byte) error
	Length(string) (int, error)
	HouseKeeping() error
}

// DataStore is used as an argument to `queue.New` to set a custom datastore
func DataStore(s Datastore) func(q *Queue) error {
	return func(q *Queue) error {
		q.store = s
		return nil
	}
}
