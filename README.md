# Q

[![CircleCI](https://circleci.com/gh/dmathieu/q/tree/master.svg?style=svg)](https://circleci.com/gh/dmathieu/q/tree/master)

A Go background worker.

## How it works

The pattern used is very similar to other background workers. You can enqueue data, which will be stored in a data store.  
Then, a second process will listen for entries pushed to the data store and execute an hanler when it gets one.

### Data Stores

Q includes two data stores by default: [memory](stores/memory.go) and [redis](stores/redis.go).  
But as long as you implement the [DataStore Interface](stores/main.go), you could implement your own with any database of your choice.

## Usage

In order to use Q, you will first need to setup a queue object.

```golang
// queue, err := q.NewQueue("memory")
queue, err := q.NewQueue("default", redisPool) // redisPool is a redigo *redis.Pool
```

You can then enqueue a job into that queue:

```golang
err := queue.Enqueue([]data("hello world"))
```

And listen for events, which needs to be done in a dedicated process

```golang
q.Run(queue, func(d []byte) error {
  log.Println(string(d))
  return nil
}, 10)
```

## "Expert" Mode

The `NewQueue` and `Run` methods are shortcuts to make the usage of Q easier.  
You may want to implement your own worker loop though.

You can then use the `q/queue` and `q/stores` packages.

Setup a queue object

```golang
// queue, err := queue.New(DataStore(&stores.MemoryStore{}))
queue, err := queue.New(queue.RedisDataStore("default", redisPool)) // redisPool is a redigo *redis.Pool
```

Enqueueing a job uses the same api as the "basic" mode

```golang
err := queue.Enqueue([]data("hello world"))
```

You can then listen for events in your own loop

```golang
for {
  err := queue.Handle(func(d []byte) error {
    log.Println(string(d))
    return nil
  })
  if err != nil {
    log.Fatal(err)
  }
}
```

Note that the `Run` method does more than just loop waiting for records.
When using this mode, you will need to handle max concurrency yourself.

You will also need to execute the `queue.HouseKeeping()` method regularly, as it recovers dead jobs.

## License

Q is released under the [MIT License](http://www.opensource.org/licenses/MIT).
