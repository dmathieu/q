package q

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisStoreIsADatastore(t *testing.T) {
	assert.Implements(t, (*Datastore)(nil), new(RedisStore))
}

func TestRedisStoreStoringAndRetrieval(t *testing.T) {
	pool := redisPool("redis://localhost:6379")
	defer pool.Close()

	t.Run("stores data", func(t *testing.T) {
		m := NewRedisStore("default", pool)
		err := m.Store([]byte("hello"))
		assert.Nil(t, err)
	})
	cleanupRedis(t, pool)

	t.Run("retrieves data", func(t *testing.T) {
		m := NewRedisStore("default", pool)
		err := m.Store([]byte("hello"))

		d, err := m.Retrieve()
		assert.Nil(t, err)
		assert.Equal(t, []byte("hello"), d)
	})
	cleanupRedis(t, pool)

	t.Run("cannot retrieve data twice", func(t *testing.T) {
		m := NewRedisStore("default", pool)
		err := m.Store([]byte("hello"))

		m.Retrieve()
		d, err := m.Retrieve()
		assert.Nil(t, err)
		assert.Nil(t, d)
	})
	cleanupRedis(t, pool)
}

func TestRedisStoreLength(t *testing.T) {
	pool := redisPool("redis://localhost:6379")
	defer pool.Close()
	m := NewRedisStore("default", pool)

	l, err := m.Length()
	assert.Nil(t, err)
	assert.Equal(t, 0, l)
	m.Store([]byte("hello"))

	l, err = m.Length()
	assert.Nil(t, err)
	assert.Equal(t, 1, l)

	cleanupRedis(t, pool)
}
