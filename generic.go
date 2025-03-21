package beacon

import "fmt"

// TypedHandler is a handler that expects a specific data type.
type TypedHandler[T any] func(TypedEvent[T]) error

// TypedEvent is an event that contains a specific data type.
type TypedEvent[T any] struct {
	Event
	Data T
}

// Wrap wraps a handler that expects a specific data type.
func Wrap[T any](handler TypedHandler[T]) Handler {
	return func(e Event) error {
		v, ok := e.Data.(T)
		if !ok {
			return fmt.Errorf("unexpected data type in wrapped handler: %T", e.Data)
		}
		return handler(TypedEvent[T]{Event: e, Data: v})
	}
}
