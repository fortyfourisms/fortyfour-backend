package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// Producer
type Producer struct {
	channel *amqp.Channel
}

// NewProducer
func NewProducer(channel *amqp.Channel) *Producer {
	return &Producer{
		channel: channel,
	}
}

// Publish
func (p *Producer) Publish(ctx context.Context, exchange, routingKey string, message interface{}) error {
	// Convert message to JSON
	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Set timeout untuk publish
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Publish message
	err = p.channel.PublishWithContext(
		ctx,
		exchange,   // exchange
		routingKey, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // Make message persistent
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}

	log.Printf("Published message to exchange=%s, routingKey=%s", exchange, routingKey)
	return nil
}
