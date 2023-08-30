package main

import (
	"context"
	"log"
	"time"

	"goq/pkg/pubsub"
	"goq/pkg/worker"
)

type Greeter struct {
	Pubsub pubsub.Pubsuber
	Name   string
}

// Handle handles the message
func (g Greeter) Handle(ctx context.Context, msg pubsub.Message) error {
	log.Println("msg: ", msg)
	return nil
}

// Delay send message to the queue
func (g Greeter) Delay(ctx context.Context, msg pubsub.Message) error {
	g.Pubsub.Publish(g.Name, msg)
	return nil
}

// GetMsgChannel returns a channel of messages
func (g Greeter) GetMsgChannel(ctx context.Context) <-chan pubsub.Message {
	return g.Pubsub.Subscribe("asdf")
}

func NewGreeter(pubsub pubsub.Pubsuber, name string) Greeter {
	return Greeter{
		Pubsub: pubsub,
		Name:   name,
	}
}

func publish(pubsuber pubsub.Pubsuber) {
	for {
		log.Println("publishing...")
		pubsuber.Publish("greeter", pubsub.Message{
			Payload: "hello",
		})

		time.Sleep(1 * time.Second)
	}
}

func main() {
	log.Println("start example")
	channel := make(chan pubsub.Message, 10)
	gochanPubsub := pubsub.NewGoChannel(channel)

	greeter := NewGreeter(gochanPubsub, "greeter")
	go publish(gochanPubsub)
	// worker
	w := worker.NewWorker()
	w.Register(greeter)
	w.Run()
	time.Sleep(10 * time.Second)
}
