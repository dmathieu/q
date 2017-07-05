package q

import (
	"fmt"
	"time"

	"github.com/dmathieu/q/queue"
	"github.com/dmathieu/q/stores"
	"github.com/garyburd/redigo/redis"
)

// NewQueue creates a new queue object with basic defaults
func NewQueue(params ...interface{}) (*queue.Queue, error) {
	var store stores.Datastore

	// Memory store
	if len(params) == 1 && params[0] == "memory" {
		store = &stores.MemoryStore{}
	}

	// Redis store
	if len(params) == 2 {
		if pool, ok := params[1].(*redis.Pool); ok {
			store = stores.NewRedisStore(params[0].(string), pool)
		}
	}

	if store == nil {
		return nil, fmt.Errorf("unknown store parameters: %q", params)
	}

	return queue.New(queue.DataStore(store))
}

// Run starts a local worker
func Run(queue *queue.Queue, handler func([]byte) error, mc int) error {
	c := make(chan struct{}, mc)
	errCh := make(chan error)

	for {
		select {
		case <-time.After(time.Minute):
			go func() {
				err := queue.HouseKeeping()
				if err != nil {
					errCh <- err
				}
			}()
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

				err := queue.Handle(handler)
				if err != nil {
					errCh <- err
				}
			}()
		}
	}
}
