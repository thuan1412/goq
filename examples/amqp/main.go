package main

import (
	"context"
	"log"
	"time"

	"goq/pkg/pubsub"
	"goq/pkg/task"
	"goq/pkg/worker"

	"github.com/google/uuid"
)

type Greeter struct {
	task.BaseTasker
	Name  string
	Queue string
}

// Handle handles the message
func (g Greeter) Handle(ctx context.Context, msg pubsub.Message) error {
	log.Println("msg: ", msg)
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

func NewGreeter(name string, queue string) *Greeter {
	return &Greeter{
		Name:  name,
		Queue: queue,
	}
}

func publish(t *Greeter) {
	ctx := context.Background()
	for {
		log.Println("publishing...")
		err := t.Async(ctx, "hello world")
		if err != nil {
			panic(err)
		}

		time.Sleep(1 * time.Second)
	}
}

func main() {
	ctx := context.Background()
	log.Println("start example")

	cfg := pubsub.AmqpConfig{
		DNS: "amqp://guest:guest@localhost:5672",
	}
	amqpPubsub, err := pubsub.NewAmqp(cfg)
	if err != nil {
		panic(err)
	}

	greeter := NewGreeter("greeter3", "greeter3")
	// worker
	w := worker.NewWorker(amqpPubsub)
	w.Register(ctx, greeter)

	go publish(greeter)
	// w.Run()
	time.Sleep(10 * time.Second)
}
