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
	Name  string
	Queue string
}

// Handle handles the message
func (g Greeter) Handle(ctx context.Context, msg pubsub.Message) error {
	log.Println("greet: ", msg)
	return nil
}

func (g Greeter) Async(ctx context.Context, text string) error {
	msg := pubsub.Message{
		Payload:  text,
		TaskName: g.Name,
		UUID:     uuid.New().String(),
	}

	return g.Delay(ctx, msg)
}

func (g Greeter) Delay(ctx context.Context, msg pubsub.Message) error {
	return g.GetPubsuber().Publish(ctx, g.Queue, msg)
}

func (g Greeter) GetName() string {
	return g.Name
}
func (g Greeter) GetQueue() string {
	return g.Queue
}

func NewGreeter(name string, queue string) *Greeter {
	return &Greeter{
		Name:  name,
		Queue: queue,
	}
}
