package task

import (
	"context"
	"goq/pkg/pubsub"
)

// Tasker is the interface that wraps the basic methods of a task.
type Tasker interface {
	Handle(ctx context.Context, msg pubsub.Message) error
	Delay(ctx context.Context, msg pubsub.Message) error
	SetPubsuber(ctx context.Context, pubsuber pubsub.Pubsuber)
}

type BaseTasker struct {
	pubsuber pubsub.Pubsuber
}

func (b *BaseTasker) Delay(ctx context.Context, msg pubsub.Message) error {
	b.pubsuber.Publish(context.Background(), "asdf", msg)
	return nil
}

func (b *BaseTasker) SetPubsuber(_ context.Context, pubsuber pubsub.Pubsuber) {
	b.pubsuber = pubsuber
}
