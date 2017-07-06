package main

import (
	"log"
	"net/url"
	"os"
	"time"

	"github.com/dmathieu/q"
	"github.com/dmathieu/q/queue"
	"github.com/dmathieu/q/stores"
	"github.com/garyburd/redigo/redis"
	"github.com/sirupsen/logrus"
)

const (
	maxConcurrency = 10
)

func main() {
	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		logrus.Fatalf("provide a REDIS_URL environment variable")
	}
	url, err := url.Parse(redisURL)
	if err != nil {
		log.Fatalf("Invalid REDIS_URL: %s", err)
	}
	pool := redisPool(url)
	queue, err := queue.New(stores.RedisDataStore("default", pool))

	logrus.Info("Listening for events")
	q.Run(queue, func(d []byte) error {
		logrus.Info(string(d))
		return nil
	}, maxConcurrency)
}

func redisPool(u *url.URL) *redis.Pool {
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
