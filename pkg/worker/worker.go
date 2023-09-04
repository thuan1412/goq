package worker

import (
	"context"
	"errors"
	"fmt"
	"goq/pkg/pubsub"
	"goq/pkg/task"
	"log"
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

func (w *Worker) Run() {
	queues := []string{}
	for _, t := range w.taskMap {
		queues = append(queues, t.GetQueue())
	}
	msgs := w.pubsuber.Subscribe(context.Background(), queues...)
	go func() {
		for msg := range msgs {
			t, ok := w.taskMap[msg.TaskName]
			if !ok {
				log.Printf("task '%s' not found\n", msg.TaskName)
				continue
			}
			err := t.Handle(context.Background(), msg)
			if err != nil {
				log.Println("error handling message: ", err)
			}
		}
	}()
}

func (w *Worker) Register(ctx context.Context, t task.Tasker) {
	t.SetPubsuber(ctx, w.pubsuber)
	if _, ok := w.taskMap[t.GetName()]; ok {
		err := errors.New(fmt.Sprintf("task '%s' already registered", t.GetName()))
		panic(err)
	}
	w.taskMap[t.GetName()] = t
}
