package pubsub

import "context"

type GoChannel struct {
	topic chan Message
}

func (g GoChannel) Publish(ctx context.Context, topicName string, msg Message) (_ error) {
	g.topic <- msg
	return nil
}

func (g GoChannel) Subscribe(ctx context.Context, topicName string) <-chan Message {
	return g.topic
}

func NewGoChannel(channel chan Message) GoChannel {
	return GoChannel{
		topic: channel,
	}
}

func (g GoChannel) Message(topicName string, message []byte) <-chan Message {
	return g.topic
}
