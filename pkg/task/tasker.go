package task

import (
	"context"
	"goq/pkg/pubsub"
)

// Tasker is the interface that wraps the basic methods of a task.
type Tasker interface {
	Handle(ctx context.Context, msg pubsub.Message) error
	Delay(ctx context.Context, payload pubsub.Message) error
	SetPubsuber(ctx context.Context, pubsuber pubsub.Pubsuber)
	GetName() string
	GetQueue() string
}

type BaseTasker struct {
	pubsuber  pubsub.Pubsuber
	topicName string
}

func (b *BaseTasker) SetPubsuber(_ context.Context, pubsuber pubsub.Pubsuber) {
	b.pubsuber = pubsuber
}

func (b *BaseTasker) GetPubsuber() pubsub.Pubsuber {
	return b.pubsuber
}
