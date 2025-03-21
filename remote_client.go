package beacon

import (
	"context"

	protoc "github.com/YONEDASH/beacon/internal/protoc"
	"github.com/bytedance/sonic"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var (
	sonicApi = sonic.ConfigFastest
)

// remoteHttpPayload represents the data structure that is sent to the remote server.
type remoteHttpPayload struct {
	EventName string `json:"event_name"`
	Event     Event  `json:"event"`
}

func grpcPostEvent(ctx context.Context, client protoc.EventServiceClient, eventName string, e Event) error {
	data, err := sonicApi.Marshal(e.Data)
	if err != nil {
		return err
	}

	req := &protoc.SubmitEventRequest{
		EventName: eventName,
		Timestamp: timestamppb.New(e.Timestamp),
		Data:      string(data),
	}

	_, err = client.SubmitEvent(ctx, req)
	return err
}
