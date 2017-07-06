package q

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/dmathieu/q/queue"
	"github.com/dmathieu/q/stores"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {

	t.Run("with no error", func(t *testing.T) {
		queue, _ := queue.New(&stores.MemoryStore{})

		var mutex = &sync.Mutex{}
		var received [][]byte
		go func() {
			err := Run(queue, func(d []byte) error {
				mutex.Lock()
				defer mutex.Unlock()

				received = append(received, d)
				return nil
			}, 1)
			assert.Nil(t, err)
		}()

		queue.Enqueue([]byte("hello"))
		queue.Enqueue([]byte("world"))

		time.Sleep(time.Millisecond)
		mutex.Lock()
		assert.Equal(t, [][]byte{[]byte("world"), []byte("hello")}, received)
		mutex.Unlock()
	})

	t.Run("with an error", func(t *testing.T) {
		queue, _ := queue.New(&stores.MemoryStore{})

		var mutex = &sync.Mutex{}
		var received [][]byte
		go func() {
			err := Run(queue, func(d []byte) error {
				fmt.Fprintf(os.Stdout, "RUNNING\n")
				mutex.Lock()
				defer mutex.Unlock()

				received = append(received, d)
				return errors.New("an error occured")
			}, 1)
			assert.Nil(t, err)
		}()

		queue.Enqueue([]byte("hello"))
		queue.Enqueue([]byte("world"))

		time.Sleep(time.Millisecond)
		mutex.Lock()
		assert.Equal(t, [][]byte{[]byte("world"), []byte("hello")}, received)
		mutex.Unlock()
	})
}
