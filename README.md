# PubSub Library Using Actor Model in Go

## Introduction

This is a Simple **PubSub (Publish-Subscribe) architecture library in Go**. The goal of this library is to provide a simple PubSub interface that can be used by any application for managing concurrent publish and subscribe tasks. What makes this library unique is that it is built on top of the **Actor Model** for concurrency. This allows for a highly concurrent system without the need for explicit locking mechanisms like **mutexes**. Instead, we leverage **Go channels** to pass messages between actors asynchronously.

In this implementation, **actors** represent different components of the PubSub system:
- **Publisher**: An actor responsible for publishing messages to topics.
- **Subscriber**: An actor responsible for subscribing to topics and consuming messages.
- **Broker**: The central actor that handles routing messages between publishers and subscribers.

### Why the Actor Model?

By using the actor model, we can achieve a clean separation of concerns for each component, where each actor has its own state and processes its own messages without worrying about shared memory or data races. The communication between actors is handled via message passing through channels, making the system simpler to reason about and scale.

## Actor Library: Current Functionality

The **Actor** library is the foundation of this PubSub system. Below is an overview of the core functionalities and methods provided by the actor library.

---

### Core Concepts

1. **Actor**:
   - An independent unit that manages its own state and processes messages from its **message queue**.
   - Each actor has:
     - A **name** for identification.
     - An optional **parent** actor to control its lifecycle.
     - A **message queue** for receiving messages.

2. **Message**:
   - A generic message passed between actors. It contains:
     - **`Content`**: The payload of the message, which can be any type.
     - **`Sender`**: The actor that sent the message, allowing the receiver to identify its origin.
     - **`Timestamp`**: The time the message was created, useful for logging and debugging.

3. **StopSignal**:
   - A message type used to stop an actor. Only the actor itself or its parent can send a `StopSignal`, ensuring controlled shutdown.

### Current Methods and Functionalities

- **NewActor(name string, parent *Actor)**: 
  - This function initializes a new actor with a name and an optional parent. If an actor has a parent, the parent actor can stop its child.
  
- **Start()**: 
  - The `Start()` method begins the actor's event loop, where it listens for incoming messages in its message queue and processes them asynchronously.
  
- **Send(target *Actor, content interface{})**: 
  - This method allows an actor to send a message to another actor. The message is appended to the target actor's message queue along with a timestamp.
  
- **Stop(target *Actor)**: 
  - This method sends a stop signal to the target actor, allowing the target actor or its parent to gracefully shut down.

- **receive(Message)**: 
  - This method is used to process general messages. It logs the message content and simulates some work using `time.Sleep`.
  
- **handleStopSignal(StopSignal)**: 
  - This method processes a stop signal. It ensures that only the actor itself or its parent can send a stop signal to halt the actor's message loop.

---

## PubSub Library: Built on Top of the Actor Model

### How the PubSub System Works

The PubSub system consists of three key actor types: **Broker**, **Publisher**, and **Subscriber**. The **Broker** acts as the central entity responsible for routing messages between publishers and subscribers. **Publishers** send messages to topics, while **Subscribers** listen for messages on specific topics.

### Components

1. **Broker**:
   - The **Broker** actor manages subscriptions and message routing.
   - It keeps track of which subscribers are subscribed to which topics in a map: `subscribers map[string][]*actor.Actor`.
   - When a **Publisher** sends a message to a topic, the **Broker** forwards that message to all relevant subscribers.
   
2. **Publisher**:
   - A **Publisher** is an actor that sends publish requests to the broker.
   - It does not send messages directly to subscribers, but rather publishes them to topics via the broker, which then distributes them.

3. **Subscriber**:
   - A **Subscriber** is an actor that subscribes to a specific topic. It sends a subscribe request to the **Broker** and starts consuming messages routed to that topic.


## Prerequisites: Installing Go

Before running this project, you need to install **Go**. Follow the instructions below to install Go on your machine.

### Step 1: Install Go

Go to the [official Go website](https://golang.org/dl/) and download the appropriate installer for your operating system. 

Alternatively, use the following commands based on your OS:

- **For Ubuntu/Linux**:
    ```bash
    sudo apt update
    sudo apt install golang-go
    ```

- **For macOS** (using Homebrew):
    ```bash
    brew install go
    ```

- **For Windows**:
    Download and install from [https://golang.org/dl/](https://golang.org/dl/).

### Step 2: Verify the Installation

To verify that Go is installed correctly, run:

```bash
go version
```

You should see output like:

```
go version go1.16.5 linux/amd64
```

### Step 3: Set Up Your Go Workspace

Make sure your Go workspace is set up correctly. Go uses the `$GOPATH` environment variable to determine where your workspace is located. You can set up your workspace by adding the following to your `.bashrc`, `.zshrc`, or corresponding shell configuration file:

```bash
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin
```

Then reload your shell:

```bash
source ~/.bashrc  # or source ~/.zshrc
```

---

## Cloning and Running the Project

### Step 1: Clone the Repository

Once Go is installed, you can clone the repository containing the PubSub library using:

```bash
git clone https://github.com/c-urly/pubsub_demo.git
cd pubsub_demo
```

### Step 2: Initialize Go Modules

Initialize the Go modules for the project by running the following command inside the project directory:

```bash
go mod tidy
```

This will download any necessary dependencies and set up the project’s environment.

### Step 3: Run the Example Code

Inside the project, you'll find a `main.go` file that demonstrates how to use the PubSub library with actors. You can run it directly using:

```bash
go run main.go
```

To check for race conditions you can run it with 

```bash
go run -race main.go
```
---

## How to Initialize and Run the PubSub Library

### Step-by-Step Instructions

1. **Initialize the Broker**
   - The **Broker** is responsible for managing subscriptions and routing messages between publishers and subscribers.
   - To create a broker, use the `NewBroker()` function, and start it by calling `Start()`.

   ```go
   broker := pubsub.NewBroker()
   broker.Start()
   ```

2. **Create Subscribers**
   - A **Subscriber** actor subscribes to a specific topic. 
   - You can create a subscriber by calling `NewSubscriber(broker, topic)` and starting the subscriber by invoking `StartConsuming()`.

   Example:
   ```go
   subscriber1 := pubsub.NewSubscriber(broker, "topic1")
   subscriber2 := pubsub.NewSubscriber(broker, "topic2")

   go subscriber1.StartConsuming() // Subscriber 1 listens to "topic1"
   go subscriber2.StartConsuming() // Subscriber 2 listens to "topic2"
   ```

3. **Create Publishers**
   - A **Publisher** actor is responsible for publishing messages to a topic via the broker.
   - You can create a publisher by calling `NewPublisher(broker)` and publishing messages by calling `Publish(topic, content)`.

   Example:
   ```go
   publisher1 := pubsub.NewPublisher(broker)
   publisher2 := pubsub.NewPublisher(broker)

   publisher1.Publish("topic1", "Message from Publisher 1 to topic1")
   publisher2.Publish("topic2", "Message from Publisher 2 to topic2")
   ```

4. **Run the PubSub System Concurrently**
   - Both publishers and subscribers run concurrently, making it a highly concurrent system.
   - You can simulate message passing between publishers and subscribers by running the `Publish` method in a loop, with a delay to simulate real-time publishing.
   
   Example:
   ```go
   go func() {
       for i := 0; i < 5; i++ {
           publisher1.Publish("topic1", fmt.Sprintf("Message %d from Publisher 1", i))
           time.Sleep(500 * time.Millisecond) // Simulate work
       }
   }()

   go func() {
       for i := 0; i < 5; i++ {
           publisher2.Publish("topic2", fmt.Sprintf("Message %d from Publisher 2", i))
           time.Sleep(700 * time.Millisecond) // Simulate work
       }
   }()
   ```

5. **Run the Application**
   - To run your application, ensure the main function gives enough time for messages to be published and consumed before the program exits.

   Example:
   ```go
   func main() {
       // Initialize broker
       broker := pubsub.NewBroker()
       broker.Start()

       // Create subscribers
       subscriber1 := pubsub.NewSubscriber(broker, "topic1")
       subscriber2 := pubsub.NewSubscriber(broker, "topic2")

       go subscriber1.StartConsuming()
       go subscriber2.StartConsuming()

       // Create publishers
       publisher1 := pubsub.NewPublisher(broker)
       publisher2 := pubsub.NewPublisher(broker)

       // Start publishing messages concurrently
       go func() {
           for i := 0; i < 5; i++ {
               publisher1.Publish("topic1", fmt.Sprintf("Message %d from Publisher 1", i))
               time.Sleep(500 * time.Millisecond) // Simulate work
           }
       }()

       go func() {
           for i := 0; i < 5; i++ {
               publisher2.Publish("topic2", fmt.Sprintf("Message %d from Publisher 2", i))
               time.Sleep(700 * time.Millisecond) // Simulate work
           }
       }()

       // Give enough time for messages to be processed before the program exits
       time.Sleep(5 * time.Second)
   }
   ```

---

## Future Extensions

Currently, the actor-based PubSub library has the following functionality:
- Asynchronous message passing between actors using channels.
- Graceful actor shutdown via stop signals.
- Basic publish-subscribe architecture with a broker handling routing of messages.

### Possible Future Enhancements:
1. **Acknowledge Mechanism**:
   - Add support for message acknowledgments where subscribers confirm the receipt of a message.
   
2. **Actor Supervision**:
   - Implement supervision strategies for actors, where parent actors can automatically restart or manage the lifecycle of their child actors in case of failure.

3. **Persistence**:
   - Add persistence to store messages and recover the state after a failure, making the PubSub system more robust.

4. **Fault Tolerance**:
   - Implement fault-tolerant mechanisms for handling message delivery guarantees (e.g., at-least-once or exactly-once delivery).

