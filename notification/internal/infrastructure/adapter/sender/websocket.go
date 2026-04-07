package sender

import (
	"chatterbox/notification/internal/domain/entity"
	"chatterbox/notification/internal/domain/valueobject"
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
		ID:     n.ID.String(),
		UserID: n.RecepientID.String(),
		Type:   n.Type.String(),
	}

	switch n.Type {
	case valueobject.NewMessageNotificationType:
		payloadData, ok := n.Payload.(entity.MessagePayload)
		if !ok {
			return errors.New("invalid payload for notification")
		}

		notification.Payload = MessagePayload{
			ID:         payloadData.ID,
			ChatID:     payloadData.ChatID,
			SenderID:   payloadData.SenderID,
			Text:       payloadData.Text,
			OccurredAt: payloadData.OccurredAt,
		}
	default:
		return errors.New("unknown notification payload")
	}

	payload, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	if err := s.hub.Send(ctx, n.RecepientID.String(), payload); err != nil {
		return err
	}

	return nil
}
