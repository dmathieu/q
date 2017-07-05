package queue

import (
	"errors"
	"testing"

	"github.com/dmathieu/q/stores"
	"github.com/stretchr/testify/assert"
)

func TestNewQueue(t *testing.T) {
	t.Run("with no options", func(t *testing.T) {
		q, err := New()

		assert.Equal(t, errors.New("no data store specified"), err)
		assert.Nil(t, q)
	})

	t.Run("with a datastore", func(t *testing.T) {
		q, err := New(DataStore(&stores.MemoryStore{}))

		assert.Nil(t, err)
		assert.NotNil(t, q)
	})
}

func TestEnqueue(t *testing.T) {
	q, err := New(DataStore(&stores.MemoryStore{}))
	assert.Nil(t, err)

	err = q.Enqueue([]byte("hello world"))
	assert.Nil(t, err)
}

func TestHandle(t *testing.T) {
	var failure []byte
	q, err := New(
		DataStore(&stores.MemoryStore{}),
		FailureHandler(func(d []byte) error {
			failure = d
			return nil
		}),
	)
	assert.Nil(t, err)

	t.Run("with no error", func(t *testing.T) {
		err = q.Enqueue([]byte("hello world"))
		assert.Nil(t, err)

		var data []byte
		err = q.Handle(func(d []byte) error {
			data = d
			return nil
		})
		assert.Nil(t, err)
		assert.Equal(t, []byte("hello world"), data)
		assert.Equal(t, []byte(nil), failure)
	})

	t.Run("with an error", func(t *testing.T) {
		err = q.Enqueue([]byte("hello world"))
		assert.Nil(t, err)

		var data []byte
		err = q.Handle(func(d []byte) error {
			data = d
			return errors.New("an error occured")
		})
		assert.Nil(t, err)
		assert.Equal(t, []byte("hello world"), data)
		assert.Equal(t, []byte("hello world"), failure)
	})

	t.Run("with no record", func(t *testing.T) {
		var data []byte
		err = q.Handle(func(d []byte) error {
			data = d
			return nil
		})
		assert.Nil(t, err)
		assert.Nil(t, data)
	})
}

func TestHouseKeeping(t *testing.T) {
	q, err := New(DataStore(&stores.MemoryStore{}))
	assert.Nil(t, err)

	err = q.HouseKeeping()
	assert.Nil(t, err)
}
