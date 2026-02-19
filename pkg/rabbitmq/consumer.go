package rabbitmq

import (
	"context"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageHandler
type MessageHandler func(ctx context.Context, body []byte) error

// Consumer
type Consumer struct {
	channel *amqp.Channel
}

// NewConsumer
func NewConsumer(channel *amqp.Channel) *Consumer {
	return &Consumer{
		channel: channel,
	}
}

// Consume
func (c *Consumer) Consume(ctx context.Context, queueName string, handler MessageHandler) error {
	// Set QoS - prefetch 1 message at a time
	err := c.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Start consuming
	msgs, err := c.channel.Consume(
		queueName, // queue
		"",        // consumer tag
		false,     // auto-ack (manual ack for reliability)
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		nil,       // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	log.Printf("Started consuming from queue: %s", queueName)

	// Process messages
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Printf("Stopping consumer for queue: %s", queueName)
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Printf("Consumer channel closed for queue: %s", queueName)
					return
				}

				// Process message
				if err := handler(ctx, msg.Body); err != nil {
					log.Printf("❌ Error processing message from %s: %v", queueName, err)
					// Nack message (requeue)
					msg.Nack(false, true)
				} else {
					// Ack message (success)
					msg.Ack(false)
					log.Printf("Message processed successfully from queue: %s", queueName)
				}
			}
		}
	}()

	return nil
}
