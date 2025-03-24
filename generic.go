package beacon

import (
	"fmt"
	"reflect"
)

// EventName returns the generic event name for a given data type.
func EventName(v any) string {
	t := reflect.TypeOf(v)
	if t.PkgPath() == "" { // PkgPath() is empty for built-in types
		return t.Name()
	}
	return fmt.Sprintf("event.%s.%s", t.PkgPath(), t.Name())
}

// AsEvent returns the generic event name and data.
// Usage: engine.Submit(beacon.AsEvent(data))
func AsEvent(data any) (string, any) {
	return EventName(data), data
}

// TypedHandler is a handler that expects a specific data type.
type TypedHandler[T any] func(T) error

// Wrap wraps a handler that expects a specific data type.
func Wrap[T any](handler TypedHandler[T]) (string, Handler) {
	var empty T
	return EventName(empty), func(e Event) error {
		value, ok := e.Data.(T)
		if !ok {
			return fmt.Errorf("unexpected data type in wrapped TypedHandler[%T]: %T", empty, e.Data)
		}
		return handler(value)
	}
}
