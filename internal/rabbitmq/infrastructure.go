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
		"main_api.ruang_lingkup.created",
		"main_api.ruang_lingkup.updated",
		"main_api.ruang_lingkup.deleted",
		"main_api.domain.created",
		"main_api.domain.updated",
		"main_api.domain.deleted",
		"main_api.kategori.created",
		"main_api.kategori.updated",
		"main_api.kategori.deleted",
		"main_api.sub_kategori.created",
		"main_api.sub_kategori.updated",
		"main_api.sub_kategori.deleted",
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

	// Declare Exchange untuk CSIRT events
	if err := rmq.DeclareExchange("csirt.events", "topic"); err != nil {
		return fmt.Errorf("failed to declare exchange: %w", err)
	}

	// Declare Queues untuk CSIRT
	csirtQueues := []string{
		"csirt.created",
		"csirt.updated",
		"csirt.deleted",
	}

	for _, queueName := range csirtQueues {
		if _, err := rmq.DeclareQueue(queueName); err != nil {
			return fmt.Errorf("failed to declare queue %s: %w", queueName, err)
		}
	}

	// Bind Queues ke Exchange dengan routing keys
	csirtBindings := map[string]string{
		"csirt.created": "csirt.created",
		"csirt.updated": "csirt.updated",
		"csirt.deleted": "csirt.deleted",
	}

	for queueName, routingKey := range csirtBindings {
		if err := rmq.BindQueue(queueName, routingKey, "csirt.events"); err != nil {
			return fmt.Errorf("failed to bind queue %s: %w", queueName, err)
		}
	}

	// Perusahaan
	if err := rmq.DeclareExchange("perusahaan.events", "topic"); err != nil {
		return err
	}

	queuePerusahaan := []string{
		"perusahaan.created",
		"perusahaan.updated",
		"perusahaan.deleted",
	}

	for _, q := range queuePerusahaan {
		if _, err := rmq.DeclareQueue(q); err != nil {
			return err
		}
		if err := rmq.BindQueue(q, q, "perusahaan.events"); err != nil {
			return err
		}
	}

	// PIC
	if err := rmq.DeclareExchange("pic.events", "topic"); err != nil {
		return err
	}

	queuePic := []string{
		"pic.created",
		"pic.updated",
		"pic.deleted",
	}

	for _, q := range queuePic {
		if _, err := rmq.DeclareQueue(q); err != nil {
			return err
		}
		if err := rmq.BindQueue(q, q, "pic.events"); err != nil {
			return err
		}
	}

	// Jabatan
	if err := rmq.DeclareExchange("jabatan.events", "topic"); err != nil {
		return err
	}

	queueJabatan := []string{
		"jabatan.created",
		"jabatan.updated",
		"jabatan.deleted",
	}

	for _, q := range queueJabatan {
		if _, err := rmq.DeclareQueue(q); err != nil {
			return err
		}
		if err := rmq.BindQueue(q, q, "jabatan.events"); err != nil {
			return err
		}
	}

	log.Println("RabbitMQ infrastructure setup completed (Users, CSIRT, Perusahaan, PIC, Jabatan)")

	return nil
}
