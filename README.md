# Beacon

Beacon is a lightweight event handling engine for Go, designed to manage event handlers in a context-aware manner. It supports both local event handling and remote event submission via gRPC.

## Features

- Context-aware event handling
- Support for remote event submission via gRPC
- Easy-to-use API for subscribing and submitting events
- Functional options for configuration
- Optional use of generics for type-safe event handling

## Usage

### Creating an Engine

To create a new event handling engine, use the `New` function:

```go
import "github.com/yonedash/beacon"

engine := beacon.New()
```

### Subscribing to Events

To subscribe to an event, use the `Subscribe` method:

```go
handler := func(event beacon.Event) error {
    // Handle the event
    return nil
}

engine.Subscribe("event_name", handler)
```

### Submitting Events

To submit an event, use the `Submit` method:

```go
err := engine.Submit("event_name", eventData)
if err != nil {
    // Handle error
}
```

### Remote Event Submission

Beacon supports submitting events to a remote server using gRPC. This is useful for distributed systems where events need to be processed by a central server.

#### Setting Up the Remote Server

First, you need to set up a gRPC server that can receive events. Use the `RegisterEventService` function to register the event service with your gRPC server:

```go
import (
    "net"
    "google.golang.org/grpc"
    "github.com/yonedash/beacon"
)

func main() {
    addr := "127.0.0.1:8941"
    lis, err := net.Listen("tcp", addr)
    if err != nil {
        log.Fatal(err)
    }

    s := grpc.NewServer()
    engine := beacon.New()
    beacon.RegisterEventService(s, engine)

    if err := s.Serve(lis); err != nil {
        log.Fatal(err)
    }
}
```

#### Submitting Events to the Remote Server

To submit events to the remote server, create a gRPC client connection and configure the engine to use it:

```go
import (
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "github.com/yonedash/beacon"
)

func main() {
    conn, err := grpc.Dial("127.0.0.1:8941", grpc.WithTransportCredentials(insecure.NewCredentials()))
    if err != nil {
        log.Fatal(err)
    }
    defer conn.Close()

    engine := beacon.New(beacon.WithRemote(conn))

    err = engine.Submit("event_name", eventData)
    if err != nil {
        // Handle error
    }
}
```

### Receiving Remote Events

To handle events received from a remote client, you need to subscribe to the events on the server side. The server will automatically call the appropriate handlers when events are received.

#### Subscribing to Events on the Server

```go
import (
    "github.com/yonedash/beacon"
)

func main() {
    engine := beacon.New()

    handler := func(event beacon.Event) error {
        // Handle the event
        return nil
    }

    engine.Subscribe("event_name", handler)

    // Set up and start the gRPC server as shown in the previous section
}
```

When a remote client submits an event, the server will deserialize the event data and call the subscribed handlers.

### Optional Use of Generics

Beacon also supports the optional use of generics for type-safe event handling. This can be useful for ensuring that event handlers receive the expected data type. However, using generics is **optional**.

#### Using Generics with Events

To use generics with events, you can wrap your handler using the `Wrap` function and `TypedEvent` type:

```go
type CustomData struct {
    Value string
}

handler := func(e beacon.TypedEvent[CustomData]) error {
    if e.Data.Value != "test" {
        t.Errorf("expected 'test', got '%s'", e.Data.Value)
    }
    return nil
}

engine := beacon.New()
engine.Subscribe(beacon.Wrap(handler))

data := CustomData{Value: "test"}
if err := engine.Submit(beacon.Typed(data)); err != nil {
    // Handle error
}
```
