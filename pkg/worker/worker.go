package worker

import (
	"context"
	"goq/pkg/task"
	"log"
)

type Worker struct {
	tasks []task.Tasker
}

func NewWorker() *Worker {
	return &Worker{
		tasks: []task.Tasker{},
	}
}

func (w *Worker) Run() {
	for _, t := range w.tasks {
		go func(t task.Tasker) {
			for msg := range t.GetMsgChannel(context.Background()) {
				log.Println("msg: ", msg)
			}
		}(t)
	}
}

func (w *Worker) Register(t task.Tasker) {
	w.tasks = append(w.tasks, t)
}
