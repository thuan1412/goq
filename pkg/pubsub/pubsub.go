package pubsub

import "context"

type Handler = func(message any) error

type Message struct {
	Payload any
}

type Pubsuber interface {
	Publish(ctx context.Context, topicName string, msg Message) error
	Subscribe(ctx context.Context, topicName string) chan Message
}
