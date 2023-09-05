package pubsub

import (
	"context"
	"encoding/json"
	"log"
)

type Handler = func(message any) error

type Message struct {
	TaskName string `json:"task_name"`
	UUID     string `json:"uuid"`
	Payload  any    `json:"payload"`
}

// Marshal returns the payload as a string in json format
func (m Message) Marshal() []byte {
	jsonStr, err := json.Marshal(m)
	if err != nil {
		log.Println("error marshalling payload: ", err)
	}
	return jsonStr
}

type Pubsuber interface {
	Publish(ctx context.Context, topicName string, msg Message) error
	Subscribe(ctx context.Context, topicNames ...string) <-chan Message
	Close(ctx context.Context) error
}
