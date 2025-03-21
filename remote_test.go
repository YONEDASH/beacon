package beacon_test

import (
	"net"
	"testing"

	"github.com/YONEDASH/beacon"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestRemote(t *testing.T) {
	addr := "127.0.0.1:8941"

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		t.Fatal(err)
	}

	s := grpc.NewServer()
	receiver := beacon.New()
	beacon.RegisterEventService(s, receiver)

	go func() {
		if err := s.Serve(lis); err != nil {
			t.Fatal(err)
		}
	}()

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatal(err)
	}
	sender := beacon.New(beacon.WithRemote(conn))

	message := ""
	handler := func(e beacon.Event) error {
		message = e.Data.(string)
		return nil
	}
	receiver.Subscribe("test", handler)

	if err := sender.Submit("test", "hello world"); err != nil {
		t.Fatal(err)
	}

	if message != "hello world" {
		t.Errorf("unexpected message: %s", message)
	}
}
