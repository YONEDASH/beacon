package beacon_test

import (
	"errors"
	"testing"

	"github.com/yonedash/beacon"
)

func TestSubscribe(t *testing.T) {
	handler := func(beacon.Event) error {
		return nil
	}

	engine := beacon.New()
	engine.Subscribe("test", handler)

	if engine.Size() != 1 {
		t.Error("handler not registered: invalid size")
	}
}

func TestSubmit(t *testing.T) {
	success := false
	handler := func(beacon.Event) error {
		success = true
		return nil
	}

	engine := beacon.New()
	engine.Subscribe("test", handler)

	if err := engine.Submit("test", nil); err != nil {
		t.Fatal(err)
	}

	if !success {
		t.Error("handler not called")
	}
}

func TestSubmitError(t *testing.T) {
	handler := func(beacon.Event) error {
		return errors.New("some error message")
	}

	engine := beacon.New()
	engine.Subscribe("test", handler)

	if err := engine.Submit("test", nil); err == nil {
		t.Error("no error received from handler")
	}
}

func TestEventCounter(t *testing.T) {
	type Counter struct {
		Count int
	}

	handler := func(e beacon.Event) error {
		e.Data.(*Counter).Count++
		return nil
	}

	ctr := new(Counter)

	engine := beacon.New()
	engine.Subscribe("increment", handler)
	engine.Submit("increment", ctr)
	engine.Submit("increment", ctr)
	engine.Submit("increment", ctr)

	if ctr.Count != 3 {
		t.Error("counter was not incremented")
	}
}

func TestSubmitIncorrectEventName(t *testing.T) {
	fail := false

	handler := func(e beacon.Event) error {
		fail = true
		return nil
	}

	engine := beacon.New()
	engine.Subscribe("test", handler)
	engine.Submit("does not exist", nil)

	if fail {
		t.Error("handler was called")
	}
}
