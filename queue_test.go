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
