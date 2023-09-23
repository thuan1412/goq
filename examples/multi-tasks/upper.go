package main

import (
	"context"
	"goq/pkg/pubsub"
	"goq/pkg/task"
	"log"
	"strings"

	"github.com/google/uuid"
)

type Upper struct {
	task.BaseTasker
	Config task.Config
}

// Handle handles the message
func (g Upper) Handle(ctx context.Context, msg pubsub.Message) error {
	str := msg.Payload.(string)
	log.Println("upper: ", strings.ToUpper(str))
	return nil
}

func (g *Upper) Async(ctx context.Context, text string) error {
	msg := pubsub.Message{
		Payload:  text,
		TaskName: task.GetTaskName(g),
		UUID:     uuid.New().String(),
	}

	return g.Delay(ctx, msg)
}

func (g *Upper) Delay(ctx context.Context, msg pubsub.Message) error {
	return g.GetPubsuber().Publish(ctx, task.GetTaskQueue(g), msg)
}

func NewUpper(name string, queue string) *Upper {
	return &Upper{
		Config: task.Config{
			Name:  name,
			Queue: queue,
		},
	}
}
