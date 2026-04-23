package rabbitmq

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// MessageHandler handles incoming messages.
type MessageHandler func(ctx context.Context, body []byte) error

// Consumer handles queue subscription and message processing.
type Consumer struct {
	rmq     *RabbitMQ     // Reference to RabbitMQ wrapper for robust channel management
	channel *amqp.Channel // Direct channel for legacy use
}

// NewConsumer creates a consumer that manages its own channels and handles reconnections.
func NewConsumer(rmq *RabbitMQ) *Consumer {
	return &Consumer{
		rmq: rmq,
	}
}

// NewConsumerWithChannel creates a basic consumer with a specific pre-opened channel.
func NewConsumerWithChannel(channel *amqp.Channel) *Consumer {
	return &Consumer{
		channel: channel,
	}
}

// Consume starts consuming messages from the specified queue.
// It uses dedicated channels per queue and automatically handles reconnections if created via NewConsumer.
func (c *Consumer) Consume(ctx context.Context, queueName string, handler MessageHandler) error {
	// If it's a channel-based consumer, use legacy logic
	if c.rmq == nil {
		return c.startChannelConsume(ctx, queueName, handler)
	}

	// Managed logic: Start a supervisor goroutine for this specific queue
	go c.supervisor(ctx, queueName, handler)
	return nil
}

// supervisor manages the lifecycle of a single queue consumer.
func (c *Consumer) supervisor(ctx context.Context, queueName string, handler MessageHandler) {
	for {
		select {
		case <-ctx.Done():
			log.Printf("Shutting down supervisor for queue: %s", queueName)
			return
		default:
			// Attempt to consume
			err := c.runConsume(ctx, queueName, handler)
			if err != nil {
				log.Printf("⚠️ Consumer for %s failed: %v. Retrying in 5 seconds...", queueName, err)

				// Wait before retry
				select {
				case <-ctx.Done():
					return
				case <-time.After(5 * time.Second):
					continue
				}
			}
		}
	}
}

// runConsume performs the actual AMQP consumption. It blocks until the channel is closed or context is cancelled.
func (c *Consumer) runConsume(ctx context.Context, queueName string, handler MessageHandler) error {
	// 1. Get a fresh dedicated channel
	ch, err := c.rmq.NewChannel()
	if err != nil {
		return fmt.Errorf("failed to open channel: %w", err)
	}
	defer ch.Close()

	// 2. Set QoS - prefetch 1 message at a time for better load balancing
	if err := ch.Qos(1, 0, false); err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// 3. Start consuming
	msgs, err := ch.Consume(
		queueName,
		"",    // consumer tag (generated)
		false, // auto-ack: false (we handle manual ack for reliability)
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		return fmt.Errorf("failed to start consuming: %w", err)
	}

	log.Printf("✅ Consumer registered and listening on queue: %s", queueName)

	// 4. Listen for channel closure
	notifyClose := ch.NotifyClose(make(chan *amqp.Error))

	for {
		select {
		case <-ctx.Done():
			return nil
		case closeErr := <-notifyClose:
			if closeErr != nil {
				return fmt.Errorf("channel closed unexpectedly: %v", closeErr)
			}
			return nil // Normal closure
		case msg, ok := <-msgs:
			if !ok {
				return fmt.Errorf("delivery channel closed")
			}

			// Process message
			if err := handler(ctx, msg.Body); err != nil {
				log.Printf("❌ Error processing message from %s: %v", queueName, err)
				// Nack message and requeue so it's not lost
				msg.Nack(false, true)
			} else {
				// Ack message on success
				msg.Ack(false)
			}
		}
	}
}

// startChannelConsume maintains support for specific channel-based consumers.
func (c *Consumer) startChannelConsume(ctx context.Context, queueName string, handler MessageHandler) error {
	if c.channel == nil {
		return fmt.Errorf("consumer has no channel")
	}

	// Set QoS
	_ = c.channel.Qos(1, 0, false)

	msgs, err := c.channel.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		return err
	}

	log.Printf("Started channel-based consumer for queue: %s", queueName)

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg, ok := <-msgs:
				if !ok {
					log.Printf("Channel-based consumer closed for queue: %s", queueName)
					return
				}
				if err := handler(ctx, msg.Body); err != nil {
					log.Printf("❌ Error in channel-based consumer %s: %v", queueName, err)
					msg.Nack(false, true)
				} else {
					msg.Ack(false)
				}
			}
		}
	}()

	return nil
}
