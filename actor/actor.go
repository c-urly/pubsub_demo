package actor

import (
	"log"
	"time"
)

// Message passed between actors
type Message struct {
	Topic     string
	Content   interface{}
	Sender    *Actor
	Timestamp time.Time
}


type StopSignal struct {
	Sender *Actor
}


type Actor struct {
	name         string
	Parent       *Actor
	messageQueue chan interface{}
}


func NewActor(name string, parent *Actor) *Actor {
	return &Actor{
		name:         name,
		Parent:       parent,
		messageQueue: make(chan interface{}, 100),
	}
}

// Actor's message queue.
func (a *Actor) MessageQueue() chan interface{} {
	return a.messageQueue
}


func (a *Actor) GetName() string {
	return a.name
}

// Start initializes the actor's event loop to process messages.
func (a *Actor) Start() {
	go func() {
		for msg := range a.messageQueue {
			switch m := msg.(type) {
			case Message:
				a.receive(m)
			case StopSignal:
				a.handleStopSignal(m)
			}
		}
	}()
}

// Processes an incoming message.
func (a *Actor) receive(msg Message) {
	log.Printf("[%s] Received message on topic '%s': %+v at %v", a.GetName(), msg.Topic, msg.Content, msg.Timestamp)
	time.Sleep(1 * time.Second) // Simulating work
	log.Printf("[%s] Finished processing message: %+v", a.GetName(), msg.Content)
}

// Handles a stop signal.
func (a *Actor) handleStopSignal(signal StopSignal) {
	if signal.Sender == a || signal.Sender == a.Parent {
		log.Printf("[%s] Received stop signal from %s, stopping actor...", a.GetName(), signal.Sender.GetName())
		close(a.messageQueue)
	} else {
		log.Printf("[%s] Ignored stop signal from %s, only the parent or the actor itself can stop it.", a.GetName(), signal.Sender.GetName())
	}
}

func (a *Actor) Send(target *Actor, message Message) {
	target.messageQueue <- message
}

// Sends a stop signal to the target actor.
func (a *Actor) Stop(target *Actor) {
	stopSignal := StopSignal{Sender: a}
	target.messageQueue <- stopSignal
}
