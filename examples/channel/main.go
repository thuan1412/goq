package main

import (
	"context"
	"log"
	"time"

	"goq/pkg/pubsub"
	"goq/pkg/task"
	"goq/pkg/worker"
)

type Greeter struct {
	task.BaseTasker
	Queue string
}

// Handle handles the message
func (g Greeter) Handle(ctx context.Context, msg pubsub.Message) error {
	log.Println("msg: ", msg)
	return nil
}

func NewGreeter(queue string) *Greeter {
	return &Greeter{
		Queue: queue,
	}
}

func publish(t task.Tasker) {
	for {
		log.Println("publishing...")
		t.Delay(context.Background(), pubsub.Message{
			Payload: "hello",
		})

		time.Sleep(1 * time.Second)
	}
}

func main() {
	ctx := context.Background()
	log.Println("start example")
	channel := make(chan pubsub.Message, 10)
	gochanPubsub := pubsub.NewGoChannel(channel)

	greeter := NewGreeter("greeter")
	// worker
	w := worker.NewWorker(gochanPubsub)
	w.Register(ctx, greeter)

	go publish(greeter)
	w.Run()
	time.Sleep(10 * time.Second)
}
