package mapper

import (
	"chatterbox/notification/internal/application/port"
	"chatterbox/notification/internal/infrastructure/controller/event"
	"encoding/json"
	"errors"
)

func MapEvent(eventType string, data []byte) (port.Event, error) {
	switch eventType {
	case "message.created":
		var e event.MessageCreated
		if err := json.Unmarshal(data, &e); err != nil {
			return nil, err
		}

		return e, nil
	default:
		return nil, errors.New("unknown event type")
	}
}
