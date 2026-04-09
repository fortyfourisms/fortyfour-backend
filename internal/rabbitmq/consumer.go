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

// ConsumeCsirtCreated
func (c *Consumer) ConsumeCsirtCreated(ctx context.Context) error {
	return c.Consume(ctx, "csirt.created", func(ctx context.Context, body []byte) error {
		var event dto_event.CsirtCreatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("CSIRT Created Event: %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyCreate("csirt", event, "")
		}

		return nil
	})
}

// ConsumeCsirtUpdated
func (c *Consumer) ConsumeCsirtUpdated(ctx context.Context) error {
	return c.Consume(ctx, "csirt.updated", func(ctx context.Context, body []byte) error {
		var event dto_event.CsirtUpdatedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("CSIRT Updated Event: %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyUpdate("csirt", event, "")
		}

		return nil
	})
}

// ConsumeCsirtDeleted
func (c *Consumer) ConsumeCsirtDeleted(ctx context.Context) error {
	return c.Consume(ctx, "csirt.deleted", func(ctx context.Context, body []byte) error {
		var event dto_event.CsirtDeletedEvent
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("CSIRT Deleted Event: %+v", event)
		if c.sseService != nil {
			c.sseService.NotifyDelete("csirt", event.ID, "")
		}

		return nil
	})
}

// consumeGenericIkasEvent handles common CRUD events dynamically
func (c *Consumer) consumeGenericIkasEvent(ctx context.Context, queueName string, resource string, action string) error {
	return c.Consume(ctx, queueName, func(ctx context.Context, body []byte) error {
		var event map[string]interface{}
		if err := json.Unmarshal(body, &event); err != nil {
			return err
		}

		log.Printf("[%s] %s Event (from RabbitMQ): %+v", resource, action, event)
		if c.sseService != nil {
			userID, _ := event["user_id"].(string)
			switch action {
			case "Created":
				c.sseService.NotifyCreate(resource, event, userID)
			case "Updated":
				c.sseService.NotifyUpdate(resource, event, userID)
			case "Deleted":
				c.sseService.NotifyDelete(resource, event, userID)
			}
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
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.ruang_lingkup.created", "ruang_lingkup", "Created")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.ruang_lingkup.updated", "ruang_lingkup", "Updated")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.ruang_lingkup.deleted", "ruang_lingkup", "Deleted")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.domain.created", "domain", "Created")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.domain.updated", "domain", "Updated")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.domain.deleted", "domain", "Deleted")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.kategori.created", "kategori", "Created")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.kategori.updated", "kategori", "Updated")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.kategori.deleted", "kategori", "Deleted")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.sub_kategori.created", "sub_kategori", "Created")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.sub_kategori.updated", "sub_kategori", "Updated")
		},
		func(ctx context.Context) error {
			return c.consumeGenericIkasEvent(ctx, "main_api.sub_kategori.deleted", "sub_kategori", "Deleted")
		},
		c.ConsumeCsirtCreated,
		c.ConsumeCsirtUpdated,
		c.ConsumeCsirtDeleted,
	}

	for _, consumer := range consumers {
		if err := consumer(ctx); err != nil {
			return err
		}
	}

	log.Println("All Main API consumers started successfully")
	return nil
}
