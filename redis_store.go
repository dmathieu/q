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

// RedisDataStore configures the queue with a redis data store
func RedisDataStore(name string, pool *redis.Pool) func(q *Queue) error {
	return DataStore(&RedisStore{name, pool})
}

func (r *RedisStore) queue() string {
	return fmt.Sprintf("q:%s:queue", r.name)
}

func (r *RedisStore) workingQueue() string {
	return fmt.Sprintf("q:%s:queue:working", r.name)
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

	d, err := redis.Bytes(conn.Do("RPOPLPUSH", r.queue(), r.workingQueue()))
	if err == redis.ErrNil {
		err = nil
	}
	return d, err
}

// Finish marks a task as finished
func (r *RedisStore) Finish(d []byte, err error) error {
	if err != nil {
		return err
	}

	conn := r.pool.Get()
	defer conn.Close()

	_, err = conn.Do("LREM", r.workingQueue(), 0, d)
	return err
}

// Length returns the number of elements in the in-memory array
func (r *RedisStore) Length() (int, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("LLEN", r.queue()))
}

// WorkingLength returns the number of elements currently being processed
func (r *RedisStore) WorkingLength() (int, error) {
	conn := r.pool.Get()
	defer conn.Close()

	return redis.Int(conn.Do("LLEN", r.workingQueue()))
}
