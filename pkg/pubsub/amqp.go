package pubsub

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/samber/lo"
)

const defaultExchange = "goq"
const defaultContentType = "text/plain"
const defaultDns = "amqp://localhost:5672/"

const consumerName = "goq"

// Amqp is a struct that implements the PubSub interface for AMQP
type Amqp struct {
	conn        *amqp.Connection
	ch          *amqp.Channel
	exchange    string
	contentType string
}

func (a Amqp) Publish(ctx context.Context, topicName string, msg Message) error {
	payloadStr := msg.Marshal()

	// TODO: add code-gen in order to remove this block
	_, err := a.ch.QueueDeclare(
		topicName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}

	a.ch.QueueBind(topicName, topicName, a.exchange, false, nil)

	return a.ch.PublishWithContext(ctx, a.exchange, topicName, false, false, amqp.Publishing{
		ContentType: a.contentType,
		Body:        payloadStr,
	})
}

func (a Amqp) Subscribe(ctx context.Context, topicNames ...string) <-chan Message {
	// test with buffered channel with size of 10
	msgs := make(chan Message, 10)

	for idx, topicName := range topicNames {
		deliveryMsgs, err := a.ch.Consume(
			topicName,
			fmt.Sprintf("%s-%s-%d", consumerName, topicName, idx),
			true,
			false,
			false,
			false,
			nil,
		)
		log.Println("subscribed to queue: ", topicName)
		if err != nil {
			panic(err)
		}

		go func() {
			for deliverMsg := range deliveryMsgs {
				var msg Message
				err := json.Unmarshal(deliverMsg.Body, &msg)
				if err != nil {
					panic(err)
				}
				msgs <- msg
			}
		}()
	}

	return msgs

}

type AmqpConfig struct {
	DNS string
	// exchange name, default is goq
	Exchange *string
	// content type, default is application/json
	ContentType *string
}

func NewAmqp(config AmqpConfig) (Amqp, error) {
	dns := lo.Ternary(config.DNS != "", config.DNS, defaultDns)
	conn, err := amqp.Dial(dns)
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	exchange := lo.TernaryF(config.Exchange != nil, func() string { return *config.Exchange }, func() string { return defaultExchange })
	contentType := lo.TernaryF(config.ContentType != nil, func() string { return *config.ContentType }, func() string { return defaultContentType })

	// create exchange
	err = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)

	return Amqp{
		conn:        conn,
		ch:          ch,
		exchange:    exchange,
		contentType: contentType,
	}, nil
}
