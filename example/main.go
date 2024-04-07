package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/diego-augusto/go-worker"
)

type job struct {
	id    int
	delay time.Duration
}

func NewJob(id, delay int) *job {
	return &job{
		id:    id,
		delay: time.Duration(delay) * time.Second,
	}
}

func (j *job) Do(ctx context.Context) error {

	// Get the executor id
	executorID := ctx.Value(keyCTX).(int)

	fmt.Printf("Doing job: %d with executor: %d\n", j.id, executorID)
	time.Sleep(j.delay)
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

	// Create jobs
	jobs := make([]worker.Doer, 0)
	delayInSec := 1
	for i := range 100 {
		jobs = append(jobs, NewJob(i, delayInSec))
	}

	// Create executors
	executors := make([]worker.Executer, 0)
	for i := range 10 {
		executors = append(executors, NewExecutor(i))
	}

	w, err := worker.New(
		worker.WithDoers(jobs...),
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
		os.Exit(0)
	}()

	<-exit

	cancel()

	fmt.Println("Exiting gracefully...")
	os.Exit(0)
}
