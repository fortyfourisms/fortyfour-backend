package rabbitmq

import (
	"fmt"
	"log"

	"fortyfour-backend/pkg/rabbitmq"
)

// SetupInfrastructure
func SetupInfrastructure(rmq *rabbitmq.RabbitMQ) error {
	// Declare Exchange untuk User events
	if err := rmq.DeclareExchange("users.events", "topic"); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare Queues
	queues := []string{
		"users.created",
		"users.updated",
		"users.deleted",
		"users.password_updated",
	}

	for _, queueName := range queues {
		if _, err := rmq.DeclareQueue(queueName); err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}
	}

	// Bind Queues ke Exchange dengan routing keys
	bindings := map[string]string{
		"users.created":          "users.created",
		"users.updated":          "users.updated",
		"users.deleted":          "users.deleted",
		"users.password_updated": "users.password_updated",
	}

	for queueName, routingKey := range bindings {
		if err := rmq.BindQueue(queueName, routingKey, "users.events"); err != nil {
			return fmt.Errorf("failed to bind queue %s: %w", queueName, err)
		}
	}

	log.Println("Users RabbitMQ infrastructure setup completed")
	return nil
}
