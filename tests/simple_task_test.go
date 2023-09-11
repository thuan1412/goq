package tests

import (
	"context"
	"fmt"
	"goq/pkg/pubsub"
	"goq/pkg/task"
	"goq/pkg/worker"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

type FileWriter struct {
	task.BaseTasker
	Name  string
	Queue string
}

const tmpFolder = "./tests/simple_task/"

type GreeterPayload struct {
	filepath string
	content  string
}

// Handle handles the message
func (g FileWriter) Handle(ctx context.Context, msg pubsub.Message) error {
	// write content to file
	payload := msg.Payload.(GreeterPayload)
	err := os.WriteFile(payload.filepath, []byte(payload.content), 0644)
	if err != nil {
		return err
	}
	return nil
}

func (g FileWriter) Async(ctx context.Context, payload GreeterPayload) error {
	msg := pubsub.Message{
		Payload:  payload,
		TaskName: g.Name,
		UUID:     uuid.New().String(),
	}
	return g.Delay(ctx, msg)
}

func (g FileWriter) Delay(ctx context.Context, msg pubsub.Message) error {
	return g.GetPubsuber().Publish(ctx, g.Queue, msg)
}

func (g FileWriter) GetName() string {
	return g.Name
}
func (g FileWriter) GetQueue() string {
	return g.Queue
}

func TestSingleTask(t *testing.T) {
	ctx := context.Background()

	// time.Sleep(30 * time.Minute)
	cfg := pubsub.AmqpConfig{
		DNS: fmt.Sprintf("amqp://guest:guest@localhost:%s/goq", mappedPort.Port()),
	}
	amqpPubsub, err := pubsub.NewAmqp(cfg)
	if err != nil {
		panic(err)
	}

	fileWriter := FileWriter{
		Name:  "default",
		Queue: "default",
	}

	// worker
	w := worker.NewWorker(amqpPubsub)
	w.Register(ctx, &fileWriter)
	// publish message first

	filePath := tmpFolder + "test.txt"
	content := "hello world"
	payload := GreeterPayload{
		filepath: filePath,
		content:  content,
	}
	err = fileWriter.Async(ctx, payload)
	if err != nil {
		panic(err)
	}

	// start worker
	w.Run(ctx)

	// wait for worker 3 second
	time.Sleep(3 * time.Second)

	w.Stop(ctx)

	// verify content inside file

	file, err := os.ReadFile(filePath)
	if err != nil {
		panic(err)
	}
	log.Println("file content: ", string(file))
}
