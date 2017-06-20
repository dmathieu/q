# Basic example

This is a simple example showing communication between a process enqueuing jobs and a worker.

## Running

First, run the worker:

> REDIS_URL=redis://localhost:6379 go run examples/basic/cmd/worker/main.go

Then, run the process enqueuing jobs:

> REDIS_URL=redis://localhost:6379 go run examples/basic/main.go

The second process should enqueue the current timestamp every second, which will then be processed and logged by the worker process.
