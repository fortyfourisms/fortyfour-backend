package rabbitmq

import (
	"context"
	"encoding/json"
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

// ConsumeIkasCreated
func (c *Consumer) ConsumeIkasCreated(ctx context.Context) error {
	return c.Consume(ctx, "ikas.created", func(ctx context.Context, body []byte) error {
		var event map[string]interface{}
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("IKAS Created Event: %+v", event)

		// TODO: Implementasi logic untuk handle event
		// Contoh: Kirim notifikasi, update analytics, dll

		return nil
	})
}

// ConsumeIkasUpdated
func (c *Consumer) ConsumeIkasUpdated(ctx context.Context) error {
	return c.Consume(ctx, "ikas.updated", func(ctx context.Context, body []byte) error {
		var event map[string]interface{}
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("IKAS Updated Event: %+v", event)

		// TODO: Implementasi logic untuk handle event

		return nil
	})
}

// ConsumeIkasDeleted
func (c *Consumer) ConsumeIkasDeleted(ctx context.Context) error {
	return c.Consume(ctx, "ikas.deleted", func(ctx context.Context, body []byte) error {
		var event map[string]interface{}
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("IKAS Deleted Event: %+v", event)

		// TODO: Implementasi logic untuk handle event

		return nil
	})
}

// ConsumeIkasImported
func (c *Consumer) ConsumeIkasImported(ctx context.Context) error {
	return c.Consume(ctx, "ikas.imported", func(ctx context.Context, body []byte) error {
		var event map[string]interface{}
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("IKAS Imported Event: %+v", event)

		// TODO: Implementasi logic untuk handle event

		return nil
	})
}

// ConsumeEmailNotifications
func (c *Consumer) ConsumeEmailNotifications(ctx context.Context) error {
	return c.Consume(ctx, "notifications.email", func(ctx context.Context, body []byte) error {
		var notification map[string]interface{}
		if err := json.Unmarshal(body, &notification); err != nil {
			return err
		}

		log.Printf("Email Notification Request: %+v", notification)

		// TODO: Implementasi logic untuk kirim email
		// Contoh: Gunakan SMTP atau service seperti SendGrid

		return nil
	})
}

// StartAllConsumers
func (c *Consumer) StartAllConsumers(ctx context.Context) error {
	consumers := []func(context.Context) error{
		c.ConsumeIkasCreated,
		c.ConsumeIkasUpdated,
		c.ConsumeIkasDeleted,
		c.ConsumeIkasImported,
		c.ConsumeEmailNotifications,
	}

	for _, consumer := range consumers {
		if err := consumer(ctx); err != nil {
			return err
		}
	}

	log.Println("All consumers started successfully")
	return nil
}
