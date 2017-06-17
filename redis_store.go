package q

import (
	"fmt"

	"github.com/garyburd/redigo/redis"
)

// A RedisStore stores all records data into redis
type RedisStore struct {
	name string
	pool *redis.Pool
}

// NewRedisStore creates a new redis store instance
func NewRedisStore(name string, pool *redis.Pool) *RedisStore {
	return &RedisStore{name, pool}
}

func (r *RedisStore) queue() string {
	return fmt.Sprintf("q:%s:queue", r.name)
}

// Store add the provided data to the in-memory array
func (r *RedisStore) Store(d []byte) error {
	conn := r.pool.Get()
	defer conn.Close()

	_, err := conn.Do("LPUSH", r.queue(), d)
	return err
}

// Retrieve pops the latest data from the in-memory array
func (r *RedisStore) Retrieve() ([]byte, error) {
	conn := r.pool.Get()
	defer conn.Close()

	d, err := redis.Bytes(conn.Do("LPOP", r.queue()))
	if err == redis.ErrNil {
		err = nil
	}
	return d, err
}

// Length returns the number of elements in the in-memory array
func (r *RedisStore) Length() (int, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("LLEN", r.queue()))
}
