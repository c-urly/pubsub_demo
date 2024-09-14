// This is used to show the pubsub demo library

package main

import (
    "log"
    "time"
    "strconv"


    "github.com/c-urly/pubsub_demo/pubsub"
)

func main() {
    // Step 1: Create a broker
    broker := pubsub.NewBroker()
    broker.Start()

    // Step 2: Create two subscribers subscribing to "topic1" and "topic2"
    subscriber1 := pubsub.NewSubscriber(broker, "topic1")
    subscriber2 := pubsub.NewSubscriber(broker, "topic2")

    // Start consuming messages concurrently
    go subscriber1.StartConsuming() // Subscriber 1 listens to "topic1"
    go subscriber2.StartConsuming() // Subscriber 2 listens to "topic2"

    // Step 3: Create two publishers to publish messages on "topic1" and "topic2"
    publisher1 := pubsub.NewPublisher(broker)
    publisher2 := pubsub.NewPublisher(broker)

    // Step 4: Publish messages concurrently
    go func() {
        for i := 0; i < 5; i++ {
            log.Printf("Publisher 1 sending message to topic1: Message %d", i)
            publisher1.Publish("topic1", "Message from Publisher 1 - "+ strconv.Itoa(i))
            time.Sleep(500 * time.Millisecond) // Simulating work
        }
    }()

    go func() {
        for i := 0; i < 5; i++ {
            log.Printf("Publisher 2 sending message to topic2: Message %d", i)
            publisher2.Publish("topic2", "Message from Publisher 2 - "+ strconv.Itoa(i))
            time.Sleep(700 * time.Millisecond) // Simulating work
        }
    }()

    // Step 5: Give time for all messages to be published and consumed
    time.Sleep(5 * time.Second)
}
