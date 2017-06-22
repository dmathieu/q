package stores

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStoreIsADatastore(t *testing.T) {
	assert.Implements(t, (*Datastore)(nil), new(MemoryStore))
}

func TestMemoryStoreStoringAndRetrieval(t *testing.T) {
	t.Run("stores data", func(t *testing.T) {
		m := &MemoryStore{}
		err := m.Store([]byte("hello"))
		assert.Nil(t, err)
	})

	t.Run("retrieves data", func(t *testing.T) {
		m := &MemoryStore{}
		err := m.Store([]byte("hello"))

		d, err := m.Retrieve()
		assert.Nil(t, err)
		assert.Equal(t, []byte("hello"), d)
	})

	t.Run("cannot retrieve data twice", func(t *testing.T) {
		m := &MemoryStore{}
		err := m.Store([]byte("hello"))

		m.Retrieve()
		d, err := m.Retrieve()
		assert.Nil(t, err)
		assert.Nil(t, d)
	})
}

func TestMemoryStoreFinish(t *testing.T) {
	m := &MemoryStore{}

	t.Run("no records", func(t *testing.T) {
		assert.Equal(t, errors.New("no working data found"), m.Finish([]byte("world")))
	})

	t.Run("a known record", func(t *testing.T) {
		m.Store([]byte("hello"))
		d, _ := m.Retrieve()

		l, err := m.Length("working")
		assert.Nil(t, err)
		assert.Equal(t, 1, l)
		assert.Nil(t, m.Finish(d))

		l, err = m.Length("working")
		assert.Nil(t, err)
		assert.Equal(t, 0, l)
	})

	t.Run("an unknown record", func(t *testing.T) {
		m.Store([]byte("hello"))
		m.Retrieve()

		assert.Equal(t, errors.New("unknown working record \"world\""), m.Finish([]byte("world")))
	})
}

func TestMemoryStoreLength(t *testing.T) {
	m := &MemoryStore{}

	l, err := m.Length("waiting")
	assert.Nil(t, err)
	assert.Equal(t, 0, l)
	m.Store([]byte("hello"))

	l, err = m.Length("waiting")
	assert.Nil(t, err)
	assert.Equal(t, 1, l)
}
