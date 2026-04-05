package sender

import (
	"chatterbox/notification/internal/domain/entity"
	"context"
	"encoding/json"
	"errors"
)

type WebSocketSender struct {
	hub Hub
}

func NewWebSocketSender(hub Hub) *WebSocketSender {
	return &WebSocketSender{
		hub: hub,
	}
}

func (s *WebSocketSender) Send(ctx context.Context, n entity.Notification) error {
	notification := Notification{
		ID:     n.ID,
		UserID: n.UserID,
		Type:   n.Type,
	}

	switch n.Type {
	case "new_message":
		data, ok := n.Data.(entity.MessageData)
		if !ok {
			return errors.New("invalid data for notification")
		}

		notification.Data = MessageData{
			ID:         data.ID,
			ChatID:     data.ChatID,
			SenderID:   data.SenderID,
			Text:       data.Text,
			OccurredAt: data.OccurredAt,
		}
	default:
		return errors.New("unknown notification data")
	}

	payload, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	if err := s.hub.Send(ctx, n.UserID, payload); err != nil {
		return err
	}

	return nil
}
