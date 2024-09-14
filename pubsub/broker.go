package pubsub

import (
    "github.com/c-urly/pubsub_demo/actor"
    "log"
    "time"
)

// Broker represents an actor that manages subscriptions and message routing.
type Broker struct {
    *actor.Actor
    subscribers map[string][]*actor.Actor
}


func NewBroker() *Broker {
    broker := &Broker{
        Actor:       actor.NewActor("Broker", nil),
        subscribers: make(map[string][]*actor.Actor),
    }

    go func() {
        for msg := range broker.Actor.MessageQueue() {
            broker.receive(msg.(actor.Message))
        }
    }()

    return broker
}

// Receive processes incoming messages, including publish and subscribe requests.
func (b *Broker) receive(msg actor.Message) {
    log.Printf("[%s] Received message on topic '%s' from %s", b.GetName(), msg.Topic, msg.Sender.GetName())

    if msg.Content == "subscribe" {
        b.handleSubscribeRequest(msg.Topic, msg.Sender)
    } else {
        b.handlePublishRequest(msg.Topic, msg.Content, msg.Sender)
    }
}

// Handles a publish request by forwarding the message to subscribers.
func (b *Broker) handlePublishRequest(topic string, content interface{}, sender *actor.Actor) {
    subscribers, exists := b.subscribers[topic]
    if !exists {
        log.Printf("[%s] No subscribers for topic '%s'", b.GetName(), topic)
        return
    }

    for _, subscriber := range subscribers {
        log.Printf("[%s] Forwarding message to subscriber %s on topic '%s'", b.GetName(), subscriber.GetName(), topic)

        message := actor.Message{
            Topic:     topic,
            Content:   content,
            Sender:    sender,
            Timestamp: time.Now(),
        }
        b.Send(subscriber, message)
    }
}

// Handles a subscribe request by registering the subscriber for a topic.
func (b *Broker) handleSubscribeRequest(topic string, sender *actor.Actor) {
    b.subscribers[topic] = append(b.subscribers[topic], sender)
    log.Printf("[%s] Subscriber %s subscribed to topic '%s'", b.GetName(), sender.GetName(), topic)
}
