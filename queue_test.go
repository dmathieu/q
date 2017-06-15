package q

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type testDataStore struct{}

func TestNewQueue(t *testing.T) {
	t.Run("with no options", func(t *testing.T) {
		q, err := New("default")

		assert.Nil(t, err)
		assert.NotNil(t, q)
	})

	t.Run("with a datastore", func(t *testing.T) {
		q, err := New("default", DataStore(&testDataStore{}))

		assert.Nil(t, err)
		assert.NotNil(t, q)
	})
}
