package main

import (
	"log"
	"net/url"
	"os"
	"sync/atomic"
	"time"

	"github.com/dmathieu/q"
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
	queue, err := q.New(q.RedisDataStore("default", pool))

	go func() {
		for {
			select {
			case <-time.After(time.Minute):
				err := queue.HouseKeeping()
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}()
	var totCount int64

	go func() {
		t := time.Tick(1 * time.Second)
		var s int64

		for {
			select {
			case <-t:
				s++
				v := atomic.AddInt64(&totCount, 0)
				logrus.Infof("After %d seconds, count is %d, %d/s", s, v, v/s)
			}
		}
	}()

	logrus.Info("Listening for events")
	queue.Run(func(d []byte) error {
		atomic.AddInt64(&totCount, 1)
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
