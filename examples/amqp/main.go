package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"goq/pkg/pubsub"
	"goq/pkg/worker"
)

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
	if len(os.Args) == 1 {
		fmt.Println("usage: go run main.go [worker|publisher]")
		return
	}
	ctx := context.Background()

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
	runType := os.Args[1]
	switch runType {
	case "worker":
		w.Run(ctx)
	case "publisher":
		go publish(greeter)
		time.Sleep(10 * time.Second)
	}
}
