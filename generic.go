package beacon

import (
	"fmt"
	"reflect"
)

// TypedHandler is a handler that expects a specific data type.
type TypedHandler[T any] func(TypedEvent[T]) error

// TypedEvent is an event that contains a specific data type.
type TypedEvent[T any] struct {
	Event
	Data T
}

func genericEventNameOf(v any) string {
	return fmt.Sprintf("type_%s", reflect.TypeOf(v).Name())
}

// Typed returns the generic event name and data.
// Usage: engine.Submit(beacon.Typed(data))
func Typed(data any) (string, any) {
	return genericEventNameOf(data), data
}

// Wrap wraps a handler that expects a specific data type.
func Wrap[T any](handler TypedHandler[T]) (string, Handler) {
	var v T
	return genericEventNameOf(v), func(e Event) error {
		v, ok := e.Data.(T)
		if !ok {
			return fmt.Errorf("unexpected data type in wrapped handler: %T", e.Data)
		}
		return handler(TypedEvent[T]{Event: e, Data: v})
	}
}
