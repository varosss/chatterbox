package producer

import (
	domainevent "chatterbox/chat/internal/domain/event"
	"chatterbox/chat/internal/infrastructure/adapter/event"
	"context"
	"encoding/json"
	"fmt"

	"github.com/streadway/amqp"
)

type RabbitMQProducer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	exchange   string
}

func NewRabbitMQProducer(rabbitMQURL, exchange string) (*RabbitMQProducer, error) {
	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, fmt.Errorf("failed to dial RabbitMQ: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	err = ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQProducer{
		connection: conn,
		channel:    ch,
		exchange:   exchange,
	}, nil
}

func (p *RabbitMQProducer) Produce(ctx context.Context, events ...domainevent.Event) error {
	for _, domainEvent := range events {
		var dto any

		switch e := domainEvent.(type) {
		case domainevent.MessageCreated:
			receivers := make([]string, len(e.ReceiverIDs))
			for i, id := range e.ReceiverIDs {
				receivers[i] = id.String()
			}

			dto = event.MessageCreated{
				MessageID:  e.MessageID.String(),
				ChatID:     e.ChatID.String(),
				SenderID:   e.SenderID.String(),
				Receivers:  receivers,
				Text:       e.Text,
				OccurredAt: e.OccurredAt(),
			}

		case domainevent.ChatCreated:
			dto = event.ChatCreated{
				ChatID:       e.ChatID.String(),
				Participants: e.ParticipantIDsAsUUIDs(),
				OccurredAt:   e.OccurredAt(),
			}

		default:
			continue
		}

		routingKey := domainEvent.Name()

		body, err := json.Marshal(dto)
		if err != nil {
			return fmt.Errorf("failed to marshal event: %w", err)
		}

		err = p.channel.Publish(
			p.exchange,
			routingKey,
			false,
			false,
			amqp.Publishing{
				ContentType:  "application/json",
				Body:         body,
				Timestamp:    domainEvent.OccurredAt(),
				DeliveryMode: amqp.Persistent,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to publish event: %w", err)
		}
	}

	return nil
}

func (p *RabbitMQProducer) Close() error {
	if p.channel != nil {
		if err := p.channel.Close(); err != nil {
			return err
		}
	}
	if p.connection != nil {
		if err := p.connection.Close(); err != nil {
			return err
		}
	}
	return nil
}
