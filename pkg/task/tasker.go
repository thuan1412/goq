package task

import (
	"context"
	"goq/pkg/pubsub"
)

// Tasker is the interface that wraps the basic methods of a task.
type Tasker interface {
	Handle(ctx context.Context, msg pubsub.Message) error
	Delay(ctx context.Context, msg pubsub.Message) error
	GetMsgChannel(ctx context.Context) <-chan pubsub.Message
}
