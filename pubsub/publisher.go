package pubsub

import (
    "github.com/c-urly/pubsub_demo/actor"
    "log"
    "time"
)

// Publisher represents an actor that sends publish requests to the broker.
type Publisher struct {
    *actor.Actor
    broker *Broker
}

// Creates a new publisher actor that communicates with the broker.
func NewPublisher(broker *Broker) *Publisher {
    return &Publisher{
        Actor:  actor.NewActor("Publisher", nil),
        broker: broker,
    }
}

// Publish sends a message to the broker with a specific topic and content.
func (p *Publisher) Publish(topic string, content interface{}) {
    log.Printf("[%s] Publishing message to topic '%s': %v", p.GetName(), topic, content)

    message := actor.Message{
        Topic:     topic,
        Content:   content,
        Sender:    p.Actor,
        Timestamp: time.Now(),
    }

    p.Send(p.broker.Actor, message)
}
