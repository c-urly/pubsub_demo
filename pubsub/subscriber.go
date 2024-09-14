package pubsub

import (
	"log"
	"time"

	"github.com/c-urly/pubsub_demo/actor"
)

// Subscriber represents an actor that subscribes to a topic and receives messages from the broker.
type Subscriber struct {
	*actor.Actor
	broker *Broker
	topic  string
}

func NewSubscriber(broker *Broker, topic string) *Subscriber {
	subscriber := &Subscriber{
		Actor:  actor.NewActor("Subscriber-"+topic, nil),
		broker: broker,
		topic:  topic,
	}

	message := actor.Message{
		Topic:     topic,
		Content:   "subscribe",
		Sender:    subscriber.Actor,
		Timestamp: time.Now(),
	}

	subscriber.Send(broker.Actor, message)

	return subscriber
}

// Starts the subscriber's event loop to listen for incoming messages.
func (s *Subscriber) StartConsuming() {
	go func() {
		for msg := range s.MessageQueue() {
			s.receive(msg.(actor.Message))
		}
	}()
}

// Receive processes messages forwarded by the broker.
func (s *Subscriber) receive(msg actor.Message) {
	log.Printf("[%s] Received message on topic '%s' at %v: %+v", s.GetName(), s.topic, msg.Timestamp, msg.Content)
}
