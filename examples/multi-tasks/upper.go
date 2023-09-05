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
	Name  string
	Queue string
}

// Handle handles the message
func (g Upper) Handle(ctx context.Context, msg pubsub.Message) error {
	str := msg.Payload.(string)
	log.Println("upper: ", strings.ToUpper(str))
	return nil
}

func (g Upper) Async(ctx context.Context, text string) error {
	msg := pubsub.Message{
		Payload:  text,
		TaskName: g.Name,
		UUID:     uuid.New().String(),
	}

	return g.Delay(ctx, msg)
}

func (g Upper) Delay(ctx context.Context, msg pubsub.Message) error {
	return g.GetPubsuber().Publish(ctx, g.Queue, msg)
}

func (g Upper) GetName() string {
	return g.Name
}
func (g Upper) GetQueue() string {
	return g.Queue
}

func NewUpper(name string, queue string) *Upper {
	return &Upper{
		Name:  name,
		Queue: queue,
	}
}
