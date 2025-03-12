package beacon

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// Event represents a data structure that is passed to event handlers.
type Event struct {
	Context   context.Context `json:"-"`
	Timestamp time.Time       `json:"timestamp"`
	Data      any             `json:"data"`
}

// newEvent creates an Event instance with the given context and data.
func newEvent(ctx context.Context, v any) Event {
	return Event{
		Context:   ctx,
		Timestamp: time.Now(),
		Data:      v,
	}
}

// Handler is a function that processes an event.
type Handler func(Event) error

// Option is a functional option for configuring a Show instance.
type Option func(*Engine)

// WithHttpRemote configures the Show instance to send events to a remote server.
// The client is used to send HTTP requests, and the URL is the remote server's endpoint.
// By configuring the remote, you cannot subscribe to events locally.
func WithHttpRemote(client *http.Client, url string) Option {
	return func(ls *Engine) {
		ls.httpClient = client
		ls.remoteUrl = url
	}
}

// New creates an instance of Show to manage event handlers.
func New(opts ...Option) *Engine {
	engine := &Engine{
		handlers: make(map[string][]Handler),
	}

	for _, opt := range opts {
		opt(engine)
	}

	return engine
}

// Engine manages event handlers that are triggered in a context-aware manner.
type Engine struct {
	handlers map[string][]Handler

	eventChan chan Event

	httpClient *http.Client
	remoteUrl  string
}

func (s *Engine) processEvents() {
	for event := range s.eventChan {
		go func(evt Event) {
			for _, handle := range s.handlers[evt.Data.(string)] {
				if err := handle(evt); err != nil {
					// TODO: Log error
				}
			}
		}(event)
	}
}

// hasRemote returns true if the remote server is enabled.
func (s *Engine) hasRemote() bool {
	return s.httpClient != nil && s.remoteUrl != ""
}

// Size returns the number of registered handlers for an event name.
func (s *Engine) Size() int {
	return len(s.handlers)
}

// Subscribe adds a handler function for a specific event name.
// Event names must be non-empty strings.
// You cannot subscribe to events when remote is enabled.
func (s *Engine) Subscribe(eventName string, handler Handler) {
	if s.hasRemote() {
		panic("cannot subscribe to events when remote is enabled")
	}

	s.handlers[eventName] = append(s.handlers[eventName], handler)
}

// Submit invokes the handler functions when an event is submitted.
func (s *Engine) Submit(eventName string, data any) error {
	return s.SubmitWithContext(context.Background(), eventName, data)
}

// SubmitWithContext executes all registered handlers for a specific event in the given context.
func (s *Engine) SubmitWithContext(ctx context.Context, eventName string, data any) error {
	if eventName == "" {
		return errors.New("event name is required")
	}

	event := newEvent(ctx, data)

	// If remote is enabled, send the event to the remote server
	if s.hasRemote() {
		return httpPostEvent(ctx, s.httpClient, s.remoteUrl, eventName, event)
	}

	for _, handle := range s.handlers[eventName] {
		select {
		case <-ctx.Done():
			return ctx.Err() // Respect cancellation or timeout
		default:
			if err := handle(event); err != nil {
				return err
			}
		}
	}

	return nil
}
