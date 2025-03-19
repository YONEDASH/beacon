package beacon

import (
	"context"
	"errors"
	"time"

	protoc "github.com/yonedash/beacon/internal/protoc"
	"google.golang.org/grpc"
)

// Event represents a data structure that is passed to event handlers.
type Event struct {
	Context   context.Context
	Timestamp time.Time
	Data      any
	canceled  *bool
}

// Cancel stops propagation of the event to further handlers.
func (e Event) Cancel() {
	if e.canceled != nil {
		*e.canceled = true
	}
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

// Option is a functional option for configuring a Engine instance.
type Option func(*Engine)

// WithRemote configures the Show instance to send events to a remote server using gRPC.
func WithRemote(conn *grpc.ClientConn) Option {
	return func(ls *Engine) {
		ls.grpcClient = protoc.NewEventServiceClient(conn)
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

	grpcClient protoc.EventServiceClient
}

// hasRemote returns true if the remote server is enabled.
func (s *Engine) hasRemote() bool {
	return s.grpcClient != nil
}

// Size returns the number of registered handlers for an event name.
func (s *Engine) Size() int {
	return len(s.handlers)
}

// Subscribe adds a handler function for a specific event name.
// Event names must be non-empty strings.
func (s *Engine) Subscribe(eventName string, handler Handler) {
	s.handlers[eventName] = append(s.handlers[eventName], handler)
}

// Submit invokes the handler functions when an event is submitted.
func (s *Engine) Submit(eventName string, data any) error {
	return s.SubmitWithContext(context.Background(), eventName, data)
}

// SubmitWithContext invokes the handler functions when an event is submitted with a context.
func (s *Engine) SubmitWithContext(ctx context.Context, eventName string, data any) error {
	if eventName == "" {
		return errors.New("event name is required")
	}

	event := newEvent(ctx, data)

	// If remote is enabled, send the event to the remote server
	if s.hasRemote() {
		if err := grpcPostEvent(ctx, s.grpcClient, eventName, event); err != nil {
			return err
		}
	}

	result := make(chan error)
	go func() {
		result <- s.fireEvent(eventName, event)
	}()

	select {
	case <-event.Context.Done():
		return event.Context.Err() // Respect cancelation or timeout
	case err := <-result:
		return err
	}
}

// fireEvent executes all registered handlers for a specific event.
func (s *Engine) fireEvent(eventName string, event Event) error {
	handlers, ok := s.handlers[eventName]
	if !ok {
		return nil
	}

	event.canceled = new(bool)

	for _, handle := range handlers {
		if err := handle(event); err != nil {
			return err
		}
		if *event.canceled {
			return nil
		}
		select {
		case <-event.Context.Done():
			return event.Context.Err()
		default:
			// Continue to the next handler if the context is not done
		}
	}
	return nil
}
