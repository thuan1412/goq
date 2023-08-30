package pubsub

type Handler = func(message any) error

type Message struct {
	Payload any
}

type Pubsuber interface {
	Publish(topicName string, msg Message) error
	Subscribe(topicName string) chan Message
}
