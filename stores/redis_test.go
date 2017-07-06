package stores

import (
	"strconv"
	"testing"

	"github.com/dmathieu/q/queue"
	"github.com/stretchr/testify/assert"
)

func TestRedisStoreIsADatastore(t *testing.T) {
	assert.Implements(t, (*queue.Datastore)(nil), new(RedisStore))
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

		l, err := m.Length("working")
		assert.Nil(t, err)
		assert.Equal(t, 0, l)

		d, err := m.Retrieve()
		assert.Nil(t, err)
		assert.Equal(t, []byte("hello"), d)

		l, err = m.Length("working")
		assert.Nil(t, err)
		assert.Equal(t, 1, l)
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

	err := m.Store([]byte("hello"))
	assert.Nil(t, err)

	d, err := m.Retrieve()
	assert.Nil(t, err)

	assert.Equal(t, err, m.Finish(d))
	l, err := m.Length("working")
	assert.Nil(t, err)
	assert.Equal(t, 0, l)
}

func TestRedisStoreLength(t *testing.T) {
	pool := redisPool("redis://localhost:6379")
	defer pool.Close()
	m := RedisStore{"default", pool}

	l, err := m.Length("waiting")
	assert.Nil(t, err)
	assert.Equal(t, 0, l)
	m.Store([]byte("hello"))

	l, err = m.Length("waiting")
	assert.Nil(t, err)
	assert.Equal(t, 1, l)

	cleanupRedis(t, pool)
}

func TestRedisHouseKeeping(t *testing.T) {
	pool := redisPool("redis://localhost:6379")
	defer pool.Close()
	m := RedisStore{"default", pool}

	t.Run("with no record", func(t *testing.T) {
		err := m.HouseKeeping()
		assert.Nil(t, err)
	})
	cleanupRedis(t, pool)

	t.Run("with enqueued records", func(t *testing.T) {
		m.Store([]byte("hello"))
		m.Retrieve()

		l, _ := m.Length("working")
		assert.Equal(t, 1, l)
		err := m.HouseKeeping()
		assert.Nil(t, err)
		l, _ = m.Length("working")
		assert.Equal(t, 1, l)
	})
	cleanupRedis(t, pool)

	t.Run("with expired records", func(t *testing.T) {
		conn := pool.Get()
		defer conn.Close()

		m.Store([]byte("hello"))
		d, _ := m.Retrieve()
		conn.Do("DEL", m.lockKey(d))

		l, _ := m.Length("waiting")
		assert.Equal(t, 0, l)
		l, _ = m.Length("working")
		assert.Equal(t, 1, l)

		err := m.HouseKeeping()
		assert.Nil(t, err)

		l, _ = m.Length("waiting")
		assert.Equal(t, 1, l)
		l, _ = m.Length("working")
		assert.Equal(t, 0, l)
	})
	cleanupRedis(t, pool)

	t.Run("with more records than the default count", func(t *testing.T) {
		conn := pool.Get()
		defer conn.Close()

		for i := 1; i <= recordsCount*2; i++ {
			m.Store([]byte(strconv.Itoa(i)))
			d, _ := m.Retrieve()
			conn.Do("DEL", m.lockKey(d))
		}

		l, _ := m.Length("waiting")
		assert.Equal(t, 0, l)
		l, _ = m.Length("working")
		assert.Equal(t, recordsCount*2, l)

		err := m.HouseKeeping()
		assert.Nil(t, err)

		l, _ = m.Length("waiting")
		assert.Equal(t, recordsCount*2, l)
		l, _ = m.Length("working")
		assert.Equal(t, 0, l)
	})
}
