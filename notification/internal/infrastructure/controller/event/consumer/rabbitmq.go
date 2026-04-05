package consumer

import (
	"context"

	"github.com/streadway/amqp"

	"chatterbox/notification/internal/application/port"
	"chatterbox/notification/internal/infrastructure/controller/event/mapper"
)

type RabbitMQConsumer struct {
	connection *amqp.Connection
	channel    *amqp.Channel
	exchange   string
	queue      string

	handlers map[string]port.EventHandler
}

func NewRabbitMQConsumer(
	rabbitMQURL string,
	exchange string,
	queueName string,
) (*RabbitMQConsumer, error) {

	conn, err := amqp.Dial(rabbitMQURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	return &RabbitMQConsumer{
		connection: conn,
		channel:    ch,
		exchange:   exchange,
		queue:      q.Name,
		handlers:   make(map[string]port.EventHandler),
	}, nil
}

func (c *RabbitMQConsumer) Register(
	eventType string,
	handler port.EventHandler,
) error {

	c.handlers[eventType] = handler

	return c.channel.QueueBind(
		c.queue,
		eventType,
		c.exchange,
		false,
		nil,
	)
}

func (c *RabbitMQConsumer) Start(ctx context.Context) error {
	msgs, err := c.channel.Consume(
		c.queue,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go func() {
		for {
			select {
			case <-ctx.Done():
				return

			case msg, ok := <-msgs:
				if !ok {
					return
				}

				c.handleMessage(ctx, msg)
			}
		}
	}()

	return nil
}

func (c *RabbitMQConsumer) handleMessage(ctx context.Context, msg amqp.Delivery) {
	eventType := msg.RoutingKey

	handler, ok := c.handlers[eventType]
	if !ok {
		_ = msg.Ack(false)
		return
	}

	event, err := mapper.MapEvent(eventType, msg.Body)
	if err != nil {
		_ = msg.Nack(false, false)
		return
	}

	if err := handler.Handle(ctx, event); err != nil {
		_ = msg.Nack(false, true)
		return
	}

	_ = msg.Ack(false)
}

func (c *RabbitMQConsumer) Close() error {
	if c.channel != nil {
		if err := c.channel.Close(); err != nil {
			return err
		}
	}
	if c.connection != nil {
		if err := c.connection.Close(); err != nil {
			return err
		}
	}
	return nil
}
