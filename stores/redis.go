package stores

import (
	"crypto/sha256"
	"fmt"

	"github.com/garyburd/redigo/redis"
)

const (
	lockDuration = 60 // seconds
)

// NewRedisStore creates a new instance of a redis store
func NewRedisStore(name string, pool *redis.Pool) *RedisStore {
	return &RedisStore{name, pool}
}

// A RedisStore stores all records data into redis
type RedisStore struct {
	name string
	pool *redis.Pool
}

func (r *RedisStore) queue() string {
	return fmt.Sprintf("q:%s:queue", r.name)
}

func (r *RedisStore) workingQueue() string {
	return fmt.Sprintf("q:%s:queue:working", r.name)
}

func (r *RedisStore) lockKey(d []byte) string {
	h := sha256.New()
	h.Write(d)

	return fmt.Sprintf("q:%s:lock:%s", r.name, h.Sum(nil))
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

	_, err = conn.Do("SETEX", r.lockKey(d), lockDuration, d)
	if err != nil {
		return nil, err
	}
	return d, err
}

// Finish marks a task as finished
func (r *RedisStore) Finish(d []byte) error {
	conn := r.pool.Get()
	defer conn.Close()

	conn.Send("MULTI")
	conn.Send("LREM", r.workingQueue(), 0, d)
	conn.Send("DEL", r.lockKey(d))
	_, err := conn.Do("EXEC")

	return err
}

// Length returns the number of elements in the in-memory array
func (r *RedisStore) Length(q string) (int, error) {
	conn := r.pool.Get()
	defer conn.Close()

	var queue string
	switch q {
	case "waiting":
		queue = r.queue()
	case "working":
		queue = r.workingQueue()
	default:
		return 0, fmt.Errorf("unknown queue %s", q)
	}

	return redis.Int(conn.Do("LLEN", queue))
}
