# Beacon

Beacon is a lightweight event handling engine for Go, designed to manage event handlers in a context-aware manner. It supports both local event handling and remote event submission via HTTP.

## Features

- Context-aware event handling
- Support for remote event submission
- Easy-to-use API for subscribing and submitting events
- Functional options for configuration

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

To configure the engine to send events to a remote server, use the `WithHttpRemote` option:

```go
client := &http.Client{}
url := "http://remote-server.com/events"

engine := beacon.New(beacon.WithHttpRemote(client, url))
```

### Receiving Remote Events

To receive events from a remote source, use the `ReceiveEventHandler` function:

```go
import (
    "net/http"
    "github.com/yonedash/beacon"
)

engine := beacon.New()

mux := http.NewServeMux()
mux.HandleFunc("/events", beacon.ReceiveEventHandler(engine))

http.ListenAndServe(":8080", mux)
```
