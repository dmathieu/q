package q

import (
	"log"
	"net/url"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/stretchr/testify/assert"
)

func redisPool(addr string) *redis.Pool {
	u, err := url.Parse(addr)
	if err != nil {
		log.Fatal(err)
	}
	auth := ""
	if u.User != nil {
		if password, ok := u.User.Password(); ok {
			auth = password
		}
	}

	return &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", u.Host)
			if err != nil {
				return nil, err
			}

			if len(auth) > 0 {
				if _, err := c.Do("AUTH", auth); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

func cleanupRedis(t *testing.T, pool *redis.Pool) error {
	conn := pool.Get()
	defer conn.Close()

	keys, err := redis.ByteSlices(conn.Do("KEYS", "q:*"))
	assert.Nil(t, err)

	for _, k := range keys {
		_, err := conn.Do("DEL", k)
		assert.Nil(t, err)
	}

	return nil
}
