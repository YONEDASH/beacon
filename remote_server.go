package beacon

import (
	"context"

	protoc "github.com/yonedash/beacon/internal/protoc"
	"google.golang.org/grpc"
)

type server struct {
	protoc.UnimplementedEventServiceServer
	engine *Engine
}

func (s *server) SubmitEvent(ctx context.Context, req *protoc.SubmitEventRequest) (*protoc.SubmitEventResponse, error) {
	var v any
	if err := sonicApi.Unmarshal([]byte(req.Data), &v); err != nil {
		return &protoc.SubmitEventResponse{Success: false}, err
	}

	event := Event{
		Context:   ctx,
		Timestamp: req.Timestamp.AsTime(),
		Data:      v,
	}

	if err := s.engine.fireEvent(req.EventName, event); err != nil {
		return &protoc.SubmitEventResponse{Success: false}, err
	}

	return &protoc.SubmitEventResponse{Success: true}, nil
}

func RegisterEventService(s *grpc.Server, engine *Engine) {
	protoc.RegisterEventServiceServer(s, &server{engine: engine})
}
