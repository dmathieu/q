package q

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEnqueue(t *testing.T) {
	q, err := New(DataStore(&MemoryStore{}))
	assert.Nil(t, err)

	err = q.Enqueue([]byte("hello world"))
	assert.Nil(t, err)
}
