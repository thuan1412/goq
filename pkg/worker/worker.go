package worker

import (
	"context"
	"fmt"
	"goq/pkg/pubsub"
	"goq/pkg/task"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Worker struct {
	pubsuber pubsub.Pubsuber
	taskMap  map[string]task.Tasker
}

func NewWorker(pubsuber pubsub.Pubsuber) *Worker {
	return &Worker{
		pubsuber: pubsuber,
		taskMap:  map[string]task.Tasker{},
	}
}

func (w *Worker) Run(ctx context.Context) {
	queues := []string{}
	for _, t := range w.taskMap {
		queues = append(queues, t.GetQueue())
	}
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	msgs := w.pubsuber.Subscribe(context.Background(), queues...)
	for {
		select {
		case msg := <-msgs:
			t, ok := w.taskMap[msg.TaskName]
			if !ok {
				log.Printf("task '%s' not found\n", msg.TaskName)
				continue
			}
			err := t.Handle(context.Background(), msg)
			if err != nil {
				log.Println("error handling message: ", err)
			}
			fmt.Println("received message: ", msgs)
		case sig := <-sigs:
			// we only listen to these two signals, it means this if block is redundant
			// but we keep it here for future change when we need to listen to more signals
			if sig == syscall.SIGINT || sig == syscall.SIGTERM {
				w.pubsuber.Close(ctx)
				// prevent the losing task when worker is down not gracefully
				time.Sleep(3 * time.Second)
				os.Exit(0)
			}
		}
	}
}

func (w *Worker) Stop(ctx context.Context) {
	w.pubsuber.Close(ctx)
}

func (w *Worker) Register(ctx context.Context, t task.Tasker) {
	t.SetPubsuber(ctx, w.pubsuber)
	if _, ok := w.taskMap[t.GetName()]; ok {
		err := fmt.Errorf("task '%s' already registered", t.GetName())
		panic(err)
	}
	w.taskMap[t.GetName()] = t
}
