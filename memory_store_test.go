package q

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemoryStoreIsADatastore(t *testing.T) {
	assert.Implements(t, (*Datastore)(nil), new(MemoryStore))
}
