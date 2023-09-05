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

	greeter := NewGreeter("greeter", "default-queue")
	upper := NewUpper("upper", "default-queue")
	// worker
	w := worker.NewWorker(amqpPubsub)
	w.Register(ctx, greeter)
	w.Register(ctx, upper)

	runType := os.Args[1]
	switch runType {
	case "worker":
		w.Run(ctx)
	case "publisher":
		for {
			log.Println("publishing...")
			err := greeter.Async(ctx, "hello world")
			if err != nil {
				panic(err)
			}

			err = upper.Async(ctx, "hello world")
			if err != nil {
				panic(err)
			}
			time.Sleep(1 * time.Second)
		}
	}
}
