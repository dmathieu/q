package q

import (
	"github.com/dmathieu/q/stores"
	"github.com/garyburd/redigo/redis"
)

// DataStore is used as an argument to `queue.New` to set a custom datastore
func DataStore(s stores.Datastore) func(q *Queue) error {
	return func(q *Queue) error {
		q.store = s
		return nil
	}
}

// RedisDataStore configures the queue with a redis data store
func RedisDataStore(name string, pool *redis.Pool) func(*Queue) error {
	return DataStore(stores.NewRedisStore(name, pool))
}
