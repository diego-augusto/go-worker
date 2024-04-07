# Worker

### Description
The Go Worker project demonstrates a simple implementation of a concurrent worker pool in Go. It utilizes goroutines and channels to efficiently execute multiple tasks concurrently.

### Features
1. Concurrent execution of jobs using a worker pool pattern.
2. Control over the number of jobs and worker executors.
4. Graceful shutdown handling.

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync/atomic"
	"time"

	"github.com/diego-augusto/go-worker"
)

type counter struct {
	times int64
}

func (c *counter) Add(ctx context.Context) error {
	atomic.AddInt64(&c.times, 1)
	return nil
}

func (c *counter) Get(ctx context.Context) (int64, error) {
	return atomic.LoadInt64(&c.times), nil
}

func NewCounter() *counter {
	return &counter{
		times: 0,
	}
}

type job struct {
	id      int
	delay   time.Duration
	counter *counter
}

func NewJob(id, delay int, counter *counter) *job {
	return &job{
		id:      id,
		delay:   time.Duration(delay) * time.Second,
		counter: counter,
	}
}

func (j *job) Do(ctx context.Context) error {

	// Get the executor id
	executorID := ctx.Value(keyCTX).(int)

	fmt.Printf("Doing job: %d with executor: %d\n", j.id, executorID)
	time.Sleep(j.delay)

	j.counter.Add(ctx)

	return nil
}

type KeyCTX string

var keyCTX KeyCTX = "keyCTX"

type executor struct {
	id int
}

func (e *executor) Execute(ctx context.Context, fn func(ctx context.Context) error) error {
	ctx = context.WithValue(ctx, keyCTX, e.id)
	return fn(ctx)
}

func NewExecutor(id int) *executor {
	return &executor{
		id: id,
	}
}

func main() {

	ctx := context.Background()

	// Create listener
	listener := NewCounter()

	// Create jobs
	jobs := make([]worker.Worker, 0)
	delayInSec := 1
	for i := range 100 {
		jobs = append(jobs, NewJob(i, delayInSec, listener))
	}

	// Create executors
	executors := make([]worker.Executer, 0)
	for i := range 10 {
		executors = append(executors, NewExecutor(i))
	}

	w, err := worker.NewPool(
		worker.WithWorkers(jobs...),
		worker.WithExecuters(executors...),
	)
	if err != nil {
		panic(err)
	}

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		if err = w.Run(ctx); err != nil {
			panic(err)
		}
		fmt.Println("Worker finished...")
		counter, _ := listener.Get(ctx)
		fmt.Printf("Executions: %d\n", counter)
		os.Exit(0)
	}()

	<-exit

	cancel()

	fmt.Println("Exiting gracefully...")
	os.Exit(0)
}

```