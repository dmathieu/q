package q

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewQueue(t *testing.T) {
	t.Run("with no options", func(t *testing.T) {
		q, err := New()

		assert.Equal(t, errors.New("no data store specified"), err)
		assert.Nil(t, q)
	})

	t.Run("with a datastore", func(t *testing.T) {
		q, err := New(DataStore(&MemoryStore{}))

		assert.Nil(t, err)
		assert.NotNil(t, q)
	})
}

func TestEnqueue(t *testing.T) {
	q, err := New(DataStore(&MemoryStore{}))
	assert.Nil(t, err)

	err = q.Enqueue([]byte("hello world"))
	assert.Nil(t, err)
}

func TestDequeue(t *testing.T) {
	q, err := New(DataStore(&MemoryStore{}))
	assert.Nil(t, err)

	t.Run("with no error", func(t *testing.T) {
		err = q.Enqueue([]byte("hello world"))
		assert.Nil(t, err)

		var data []byte
		err = q.Dequeue(func(d []byte) error {
			data = d
			return nil
		})
		assert.Nil(t, err)
		assert.Equal(t, []byte("hello world"), data)
	})

	t.Run("with an error", func(t *testing.T) {
		err = q.Enqueue([]byte("hello world"))
		assert.Nil(t, err)

		var data []byte
		err = q.Dequeue(func(d []byte) error {
			data = d
			return errors.New("an error occured")
		})
		assert.Equal(t, errors.New("an error occured"), err)
		assert.Equal(t, []byte("hello world"), data)
	})

	t.Run("with no record", func(t *testing.T) {
		var data []byte
		err = q.Dequeue(func(d []byte) error {
			data = d
			return nil
		})
		assert.Nil(t, err)
		assert.Nil(t, data)
	})
}
