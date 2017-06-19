package q

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRedisStoreIsADatastore(t *testing.T) {
	assert.Implements(t, (*Datastore)(nil), new(RedisStore))
}

func TestInitRedisDataStore(t *testing.T) {
	pool := redisPool("redis://localhost:6379")
	defer pool.Close()

	q, err := New(RedisDataStore("default", pool))
	assert.Nil(t, err)
	assert.NotNil(t, q)
}

func TestRedisStoreStoringAndRetrieval(t *testing.T) {
	pool := redisPool("redis://localhost:6379")
	defer pool.Close()

	t.Run("stores data", func(t *testing.T) {
		m := RedisStore{"default", pool}
		err := m.Store([]byte("hello"))
		assert.Nil(t, err)
	})
	cleanupRedis(t, pool)

	t.Run("retrieves data", func(t *testing.T) {
		m := RedisStore{"default", pool}
		err := m.Store([]byte("hello"))

		d, err := m.Retrieve()
		assert.Nil(t, err)
		assert.Equal(t, []byte("hello"), d)
	})
	cleanupRedis(t, pool)

	t.Run("cannot retrieve data twice", func(t *testing.T) {
		m := RedisStore{"default", pool}
		err := m.Store([]byte("hello"))

		m.Retrieve()
		d, err := m.Retrieve()
		assert.Nil(t, err)
		assert.Nil(t, d)
	})
	cleanupRedis(t, pool)
}

func TestRedisStoreFinish(t *testing.T) {
	pool := redisPool("redis://localhost:6379")
	defer pool.Close()
	m := RedisStore{"default", pool}

	assert.Nil(t, m.Finish(nil))

	err := errors.New("test error")
	assert.Equal(t, err, m.Finish(err))
}

func TestRedisStoreLength(t *testing.T) {
	pool := redisPool("redis://localhost:6379")
	defer pool.Close()
	m := RedisStore{"default", pool}

	l, err := m.Length()
	assert.Nil(t, err)
	assert.Equal(t, 0, l)
	m.Store([]byte("hello"))

	l, err = m.Length()
	assert.Nil(t, err)
	assert.Equal(t, 1, l)

	cleanupRedis(t, pool)
}
