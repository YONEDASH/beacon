package beacon_test

import (
	"context"
	"errors"
	"testing"
	"time"

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

func TestEventCancel(t *testing.T) {
	counter := 0

	engine := beacon.New()
	engine.Subscribe("test", func(e beacon.Event) error {
		counter++
		e.Cancel()
		return nil
	})
	engine.Subscribe("test", func(e beacon.Event) error {
		counter++
		return nil
	})
	engine.Submit("test", nil)

	if counter == 2 {
		t.Error("handler was not cancelled and second handler was called")
	} else if counter != 1 {
		t.Error("no handler was called")
	}
}

func TestContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel()

	engine := beacon.New()
	engine.Subscribe("test", func(e beacon.Event) error {
		<-e.Context.Done()
		return nil
	})

	if err := engine.SubmitWithContext(ctx, "test", nil); err != context.DeadlineExceeded {
		t.Error("context error not received")
	}
}
