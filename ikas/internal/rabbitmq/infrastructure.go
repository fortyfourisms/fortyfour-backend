package rabbitmq

import (
	"fmt"
	"log"

	"fortyfour-backend/pkg/rabbitmq"
)

// SetupInfrastructure
func SetupInfrastructure(rmq *rabbitmq.RabbitMQ) error {
	// Declare Exchange untuk IKAS events
	if err := rmq.DeclareExchange("ikas.events", "topic"); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare Queues
	queues := []string{
		"ikas.created",
		"ikas.updated",
		"ikas.deleted",
		"ikas.imported",
		"notifications.email",
	}

	for _, queueName := range queues {
		if _, err := rmq.DeclareQueue(queueName); err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}
	}

	// Bind Queues ke Exchange dengan routing keys
	bindings := map[string]string{
		"ikas.created":        "ikas.created",
		"ikas.updated":        "ikas.updated",
		"ikas.deleted":        "ikas.deleted",
		"ikas.imported":       "ikas.imported",
		"notifications.email": "notification.email",
	}

	for queueName, routingKey := range bindings {
		if err := rmq.BindQueue(queueName, routingKey, "ikas.events"); err != nil {
			return fmt.Errorf("failed to bind queue %s: %w", queueName, err)
		}
	}

	log.Println("RabbitMQ infrastructure setup completed")
	return nil
}
