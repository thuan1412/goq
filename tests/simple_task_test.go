package tests

import (
	"context"
	"fmt"
	"goq/pkg/pubsub"
	"goq/pkg/task"
	"goq/pkg/worker"
	"io/fs"
	"log"
	"os"
	"testing"
	"time"

	"github.com/google/uuid"
)

type FileWriter struct {
	task.BaseTasker
	Config task.Config
}

const tmpFolder = "./tests/simple_task/"

type GreeterPayload struct {
	Filepath string `json:"filepath,omitempty"`
	Content  string `json:"content,omitempty"`
}

// Handle handles the message
func (g FileWriter) Handle(ctx context.Context, msg pubsub.Message) error {
	// write content to file
	payloadMap := msg.Payload.(map[string]interface{})
	payload := GreeterPayload{
		Filepath: payloadMap["filepath"].(string),
		Content:  payloadMap["content"].(string),
	}
	var term fs.FileMode = 0644
	err := os.MkdirAll(tmpFolder, term)
	if err != nil {
		fmt.Println("error creating dir:", err)
		return err
	}
	err = os.WriteFile(payload.Filepath, []byte(payload.Content), term)
	if err != nil {
		fmt.Println("error writing file:", err)
		return err
	}
	return nil
}

func (g *FileWriter) Async(ctx context.Context, payload GreeterPayload) error {
	msg := pubsub.Message{
		Payload:  payload,
		TaskName: task.GetTaskName(g),
		UUID:     uuid.New().String(),
	}
	return g.Delay(ctx, msg)
}

func (g *FileWriter) Delay(ctx context.Context, msg pubsub.Message) error {
	return g.GetPubsuber().Publish(ctx, task.GetTaskQueue(g), msg)
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

	fileWriter := FileWriter{}

	// worker
	w := worker.NewWorker(amqpPubsub)
	w.Register(ctx, &fileWriter)
	// publish message first

	filePath := tmpFolder + "test.txt"
	content := "hello world"
	payload := GreeterPayload{
		Filepath: filePath,
		Content:  content,
	}
	err = fileWriter.Async(ctx, payload)
	if err != nil {
		panic(err)
	}

	// start worker
	go w.Run(ctx)

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
