package rabbitmq

import (
	"fmt"
	"log"
	"strings"

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

	// Declare and Bind Queues for IKAS Events (for SSE)
	ikasQueues := []string{
		"main_api.ikas.created",
		"main_api.ikas.updated",
		"main_api.ikas.deleted",
	}

	for _, queueName := range ikasQueues {
		if _, err := rmq.DeclareQueue(queueName); err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}

		// Bind to ikas.events exchange (from IKAS service)
		routingKey := strings.TrimPrefix(queueName, "main_api.")
		if err := rmq.BindQueue(queueName, routingKey, "ikas.events"); err != nil {
			return fmt.Errorf("failed to bind queue %s to ikas.events: %w", queueName, err)
		}
	}

	log.Println("Users RabbitMQ infrastructure setup completed")
	return nil
}
