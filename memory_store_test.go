package q

import (
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

func TestMemoryStoreLength(t *testing.T) {
	m := &MemoryStore{}

	l, err := m.Length()
	assert.Nil(t, err)
	assert.Equal(t, 0, l)
	m.Store([]byte("hello"))

	l, err = m.Length()
	assert.Nil(t, err)
	assert.Equal(t, 1, l)
}
