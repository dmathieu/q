package q

import "fmt"

// A Queue allows enqueuing and listening to events
type Queue struct {
	name string
}

// New initializes a new queue
func New(name string) *Queue {
	return &Queue{name}
}

func (q *Queue) workingQueue() string {
	return fmt.Sprintf("%s:working", q.name)
}
