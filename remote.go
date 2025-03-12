package beacon

import (
	"bytes"
	"context"
	"fmt"
	"net/http"

	"github.com/bytedance/sonic"
)

var (
	sonicApi = sonic.ConfigFastest
)

// remoteHttpPayload represents the data structure that is sent to the remote server.
type remoteHttpPayload struct {
	EventName string `json:"event_name"`
	Event     Event  `json:"event"`
}

// httpPostEvent sends an event to a remote server using the given client and URL.
func httpPostEvent(ctx context.Context, client *http.Client, url, eventName string, e Event) error {
	data, err := sonicApi.Marshal(remoteHttpPayload{
		EventName: eventName,
		Event:     e,
	})
	if err != nil {
		return fmt.Errorf("could not marshal event: %v", err)
	}

	body := bytes.NewReader(data)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return fmt.Errorf("could not create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status: %s", resp.Status)
	}

	return err
}

// ReceiveEventHandler returns an http.HandlerFunc that can be used to receive events from a remote source.
func ReceiveEventHandler(shows ...*Engine) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var payload remoteHttpPayload
		if err := sonicApi.NewDecoder(r.Body).Decode(&payload); err != nil {
			http.Error(w, "could not decode payload", http.StatusBadRequest)
			return
		}

		for _, show := range shows {
			if err := show.Submit(payload.EventName, payload.Event.Data); err != nil {
				http.Error(w, "could not submit event", http.StatusInternalServerError)
				return
			}
		}

		w.WriteHeader(http.StatusOK)
	}
}
