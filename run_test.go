package q

import (
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/dmathieu/q/stores"
	"github.com/stretchr/testify/assert"
)

func TestRun(t *testing.T) {
	q, _ := New(DataStore(&stores.MemoryStore{}))

	t.Run("with no error", func(t *testing.T) {
		var mutex = &sync.Mutex{}
		var received [][]byte
		go func() {
			err := q.Run(func(d []byte) error {
				mutex.Lock()
				defer mutex.Unlock()

				received = append(received, d)
				return nil
			}, 1)
			assert.Nil(t, err)
		}()

		q.Enqueue([]byte("hello"))
		q.Enqueue([]byte("world"))

		time.Sleep(time.Millisecond)
		mutex.Lock()
		assert.Equal(t, [][]byte{[]byte("world"), []byte("hello")}, received)
		mutex.Unlock()
	})

	t.Run("with an error", func(t *testing.T) {
		var mutex = &sync.Mutex{}
		var received [][]byte
		go func() {
			err := q.Run(func(d []byte) error {
				mutex.Lock()
				defer mutex.Unlock()

				received = append(received, d)
				return errors.New("an error occured")
			}, 1)
			assert.Nil(t, err)
		}()

		q.Enqueue([]byte("hello"))
		q.Enqueue([]byte("world"))

		time.Sleep(time.Millisecond)
		mutex.Lock()
		assert.Equal(t, [][]byte{[]byte("world"), []byte("hello")}, received)
		mutex.Unlock()
	})
}
