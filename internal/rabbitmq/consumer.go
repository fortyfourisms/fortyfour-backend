package rabbitmq

import (
	"context"
	"encoding/json"
	"log"

	"fortyfour-backend/internal/dto/dto_event"
	"fortyfour-backend/pkg/rabbitmq"
)

// SSEBroadcaster defines the interface for SSE notifications to avoid import cycles
type SSEBroadcaster interface {
	NotifyCreate(resource string, data interface{}, userID string)
	NotifyUpdate(resource string, data interface{}, userID string)
	NotifyDelete(resource string, id interface{}, userID string)
}

// Consumer wrapper
type Consumer struct {
	*rabbitmq.Consumer
	sseService SSEBroadcaster
}

// NewConsumer
func NewConsumer(c *rabbitmq.Consumer, sseService SSEBroadcaster) *Consumer {
	return &Consumer{
		Consumer:   c,
		sseService: sseService,
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
		if c.sseService != nil {
			c.sseService.NotifyCreate("users", event, "")
		}

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
		if c.sseService != nil {
			c.sseService.NotifyUpdate("users", event, "")
		}

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
		if c.sseService != nil {
			c.sseService.NotifyDelete("users", event.ID, "")
		}

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
		if c.sseService != nil {
			c.sseService.NotifyUpdate("password", event, "")
		}

		return nil
	})
}

// ConsumeIkasCreated
func (c *Consumer) ConsumeIkasCreated(ctx context.Context) error {
	return c.Consume(ctx, "main_api.ikas.created", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("IKAS Created Event (from RabbitMQ): %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyCreate("ikas", event, event.UserID)
		}

		return nil
	})
}

// ConsumeIkasUpdated
func (c *Consumer) ConsumeIkasUpdated(ctx context.Context) error {
	return c.Consume(ctx, "main_api.ikas.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("IKAS Updated Event (from RabbitMQ): %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyUpdate("ikas", event, event.UserID)
		}

		return nil
	})
}

// ConsumeIkasDeleted
func (c *Consumer) ConsumeIkasDeleted(ctx context.Context) error {
	return c.Consume(ctx, "main_api.ikas.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.IkasDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("IKAS Deleted Event (from RabbitMQ): %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyDelete("ikas", event.IkasID, event.UserID)
		}

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
		c.ConsumeIkasCreated,
		c.ConsumeIkasUpdated,
		c.ConsumeIkasDeleted,
	}

	for _, consumer := range consumers {
		if err := consumer(ctx); err != nil {
			return err
		}
	}

	log.Println("All Main API consumers started successfully")
	return nil
}
