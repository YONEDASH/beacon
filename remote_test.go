package beacon_test

import (
	"net/http"
	"testing"

	"github.com/yonedash/beacon"
)

func TestHttpRemote(t *testing.T) {
	listenProtocol := "http://"
	listenAddr := "127.0.0.1:8941"
	listenPath := "/remote-test"

	client := http.DefaultClient

	showSender := beacon.New(beacon.WithHttpRemote(client, listenProtocol+listenAddr+listenPath))
	showListener := beacon.New()

	// Start http server
	mux := http.NewServeMux()
	mux.HandleFunc(listenPath, beacon.ReceiveEventHandler(showListener))

	go http.ListenAndServe(listenAddr, mux)

	message := ""
	handler := func(e beacon.Event) error {
		message = e.Data.(string)
		return nil
	}
	showListener.Subscribe("test", handler)

	if err := showSender.Submit("test", "hello world"); err != nil {
		t.Fatal(err)
	}

	if message != "hello world" {
		t.Errorf("unexpected message: %s", message)
	}
}
