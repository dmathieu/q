package q

import "github.com/garyburd/redigo/redis"

// Calling NewQueue with a single argument as string
// will initiate it as an in-memory store
func ExampleNewQueue_memory() {
	NewQueue("memory")
}

// Calling NewQueue with two arguments, the first one being a string
// and the second one being a *redis.Pool from redigo
// will initiate it as a redis store
func ExampleNewQueue_redis() {
	redisPool := &redis.Pool{}
	NewQueue("default", redisPool)
}
