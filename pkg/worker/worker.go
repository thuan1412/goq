package worker

import (
	"context"
	"goq/pkg/pubsub"
	"goq/pkg/task"
)

type Worker struct {
	tasks    []task.Tasker
	pubsuber pubsub.Pubsuber
}

func NewWorker(pubsuber pubsub.Pubsuber) *Worker {
	return &Worker{
		tasks:    []task.Tasker{},
		pubsuber: pubsuber,
	}
}

func (w *Worker) Run() {
	for _, t := range w.tasks {
		go func(t task.Tasker) {
			for {
				msg := <-w.pubsuber.Subscribe(context.Background(), "asdf")
				t.Handle(context.Background(), msg)
			}
		}(t)
	}
}

func (w *Worker) Register(ctx context.Context, t task.Tasker) {
	t.SetPubsuber(ctx, w.pubsuber)
	w.tasks = append(w.tasks, t)
}
