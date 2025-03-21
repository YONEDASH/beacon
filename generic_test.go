package beacon_test

import (
	"testing"

	"github.com/YONEDASH/beacon"
)

func TestWrappedEvent(t *testing.T) {
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
	engine.Subscribe("custom_event", beacon.Wrap(handler))

	data := CustomData{Value: "test"}
	if err := engine.Submit("custom_event", data); err != nil {
		t.Fatal(err)
	}
}

func TestWrappedEventIncorrectType(t *testing.T) {
	type CustomData struct {
		Value string
	}

	handler := func(e beacon.TypedEvent[CustomData]) error {
		return nil
	}

	engine := beacon.New()
	engine.Subscribe("custom_event", beacon.Wrap(handler))

	// Submitting data of incorrect type
	if err := engine.Submit("custom_event", "incorrect type"); err == nil {
		t.Error("expected error due to incorrect data type, got nil")
	}
}
