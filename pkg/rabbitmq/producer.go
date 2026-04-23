package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Producer handles message publishing.
type Producer struct {
	rmq     *RabbitMQ
	channel *amqp.Channel
}

// NewProducer creates a producer that manages its own connection state via RabbitMQ wrapper.
func NewProducer(rmq *RabbitMQ) *Producer {
	return &Producer{
		rmq: rmq,
	}
}

// NewProducerWithChannel creates a producer with a specific pre-opened channel.
func NewProducerWithChannel(channel *amqp.Channel) *Producer {
	return &Producer{
		channel: channel,
	}
}

// Publish sends a message to the specified exchange and routing key.
func (p *Producer) Publish(ctx context.Context, exchange, routingKey string, message interface{}) error {
	// Convert message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// 1. Get the channel to use
	ch := p.channel
	if p.rmq != nil {
		// Get current shared channel from RabbitMQ wrapper
		ch = p.rmq.GetChannel()
	}

	if ch == nil || ch.IsClosed() {
		return fmt.Errorf("RabbitMQ channel is closed")
	}

	// 2. Set timeout for publish operation
	publishCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 3. Publish message with persistence
	err = ch.PublishWithContext(
		publishCtx,
		exchange,
		routingKey,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message to %s -> %s", exchange, routingKey)
	return nil
}
