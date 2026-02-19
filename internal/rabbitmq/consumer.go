package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"fortyfour-backend/internal/dto/dto_event"
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

// ConsumeUserCreated
func (c *Consumer) ConsumeUserCreated(ctx context.Context) error {
	return c.Consume(ctx, "users.created", func(ctx context.Context, body []byte) error {
		var event dto_event.UserCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("USER Created Event: %+v", event)
		// TODO: Implement logic

		return nil
	})
}

// ConsumeUserUpdated
func (c *Consumer) ConsumeUserUpdated(ctx context.Context) error {
	return c.Consume(ctx, "users.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.UserUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("USER Updated Event: %+v", event)
		// TODO: Implement logic

		return nil
	})
}

// ConsumeUserDeleted
func (c *Consumer) ConsumeUserDeleted(ctx context.Context) error {
	return c.Consume(ctx, "users.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.UserDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("USER Deleted Event: %+v", event)
		// TODO: Implement logic

		return nil
	})
}

// ConsumeUserPasswordUpdated
func (c *Consumer) ConsumeUserPasswordUpdated(ctx context.Context) error {
	return c.Consume(ctx, "users.password_updated", func(ctx context.Context, body []byte) error {
		var event dto_event.UserPasswordUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("USER Password Updated Event: %+v", event)
		// TODO: Implement logic

		return nil
	})
}

// StartAllConsumers
func (c *Consumer) StartAllConsumers(ctx context.Context) error {
	consumers := []func(context.Context) error{
		c.ConsumeUserCreated,
		c.ConsumeUserUpdated,
		c.ConsumeUserDeleted,
		c.ConsumeUserPasswordUpdated,
	}

	for _, consumer := range consumers {
		if err := consumer(ctx); err != nil {
			return err
		}
	}

	log.Println("All User consumers started successfully")
	return nil
}
