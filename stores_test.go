package q

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
	"github.com/stretchr/testify/assert"
)

func TestInitRedisDataStore(t *testing.T) {
	mock := redigomock.NewConn()
	pool := mockPool(mock)

	q, err := New(RedisDataStore("default", pool))
	assert.Nil(t, err)
	assert.NotNil(t, q)
}

func mockPool(conn redis.Conn) *redis.Pool {
	return redis.NewPool(func() (redis.Conn, error) {
		return conn, nil
	}, 10)
}
