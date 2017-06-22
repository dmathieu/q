package stores

// A Datastore allows communicating with a data storage for storing and retrieving record
type Datastore interface {
	Store([]byte) error
	Retrieve() ([]byte, error)
	Finish([]byte) error
	Length(string) (int, error)
}
