package task

import (
	"context"
	"goq/pkg/pubsub"
	"testing"

	"github.com/stretchr/testify/assert"
)

type task struct {
	Config Config
}

func (*task) Handle(ctx context.Context, msg pubsub.Message) (_ error) {
	panic("not implemented") // TODO: Implement
}
func (*task) Delay(ctx context.Context, payload pubsub.Message) (_ error) {
	panic("not implemented") // TODO: Implement
}
func (*task) SetPubsuber(ctx context.Context, pubsuber pubsub.Pubsuber) {
	panic("not implemented") // TODO: Implement
}

func TestGetTaskName(t *testing.T) {
	t.Run("should return task name when specified", func(t *testing.T) {
		task := &task{
			Config: Config{
				Name: "task name",
			},
		}
		got := GetTaskName(task)
		assert.Equal(t, "task name", got)
	})

	t.Run("should return queue when specified", func(t *testing.T) {
		task := &task{
			Config: Config{
				Queue: "priority",
			},
		}
		got := GetTaskQueue(task)
		assert.Equal(t, "priority", got)
	})

	t.Run("should return default when no config provided", func(t *testing.T) {
		task := &task{}
		gotName := GetTaskName(task)
		assert.Equal(t, "*task.task", gotName)

		gotQueue := GetTaskQueue(task)
		assert.Equal(t, "default", gotQueue)
	})
}
