package main

import (
	"context"
	"goq/pkg/pubsub"
	"goq/pkg/task"
	"log"

	"github.com/google/uuid"
)

type Greeter struct {
	task.BaseTasker
	Config task.Config
}

// Handle handles the message
func (g *Greeter) Handle(ctx context.Context, msg pubsub.Message) error {
	log.Println("hello: ", msg)
	return nil
}

// Async calles Delay to publish message into queue
// TODO: remove this or Delay to simplify the interface
func (g *Greeter) Async(ctx context.Context, text string) error {
	msg := pubsub.Message{
		Payload:  text,
		TaskName: task.GetTaskName(g),
		UUID:     uuid.New().String(),
	}
	return g.Delay(ctx, msg)
}

func (g *Greeter) Delay(ctx context.Context, msg pubsub.Message) error {
	return g.GetPubsuber().Publish(ctx, task.GetTaskName(g), msg)
}

func NewGreeter(name string, queue string) *Greeter {
	return &Greeter{
		Config: task.Config{
			Name:  name,
			Queue: queue,
		},
	}
}
