package q

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func TestNewQueue(t *testing.T) {

	t.Run("with a memory store", func(t *testing.T) {
		queue, err := NewQueue("memory")
		assert.Nil(t, err)
		assert.NotNil(t, queue)
	})

	t.Run("with a redis store", func(t *testing.T) {
		pool := &redis.Pool{}
		queue, err := NewQueue("default", pool)
		assert.Nil(t, err)
		assert.NotNil(t, queue)
	})

	t.Run("with an unknown store", func(t *testing.T) {
		queue, err := NewQueue("something")
		assert.Equal(t, errors.New(`unknown store parameters: ["something"]`), err)
		assert.Nil(t, queue)
	})
}

func TestRun(t *testing.T) {

	t.Run("with no error", func(t *testing.T) {
		queue, _ := NewQueue("memory")

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
		queue, _ := NewQueue("memory")

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
