package q

import (
	"fmt"
	"time"

	"github.com/dmathieu/q/queue"
)

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
