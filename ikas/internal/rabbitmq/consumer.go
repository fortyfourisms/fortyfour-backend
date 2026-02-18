package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"fortyfour-backend/pkg/rabbitmq"
)

// Consumer
type Consumer struct {
	*rabbitmq.Consumer
}

// NewConsumer
func NewConsumer(c *rabbitmq.Consumer) *Consumer {
	return &Consumer{
		Consumer: c,
	}
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
