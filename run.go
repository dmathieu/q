package q

import (
	"fmt"
	"time"
)

// Run starts a local worker
func (q *Queue) Run(handler func([]byte) error, mc int) error {
	c := make(chan struct{}, mc)
	errCh := make(chan error)

	for {
		select {
		case <-time.After(time.Minute):
			go func() {
				err := q.HouseKeeping()
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

				err := q.Handle(handler)
				if err != nil {
					errCh <- err
				}
			}()
		}
	}
}
