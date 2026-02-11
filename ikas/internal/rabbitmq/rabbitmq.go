package rabbitmq

import (
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQ
type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	url     string
}

// NewRabbitMQ
func NewRabbitMQ(url string) (*RabbitMQ, error) {
	rmq := &RabbitMQ{
		url: url,
	}

	if err := rmq.connect(); err != nil {
		return nil, err
	}

	return rmq, nil
}

// connect
func (r *RabbitMQ) connect() error {
	var err error

	// create connection
	r.conn, err = amqp.Dial(r.url)
	if err != nil {
		return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	// create channel
	r.channel, err = r.conn.Channel()
	if err != nil {
		r.conn.Close()
		return fmt.Errorf("failed to open channel: %w", err)
	}

	log.Println("Connected to RabbitMQ successfully")
	return nil
}

// GetChannel
func (r *RabbitMQ) GetChannel() *amqp.Channel {
	return r.channel
}

// Close
func (r *RabbitMQ) Close() {
	if r.channel != nil {
		r.channel.Close()
	}
	if r.conn != nil {
		r.conn.Close()
	}
	log.Println("🔌 RabbitMQ connection closed")
}

// DeclareExchange
func (r *RabbitMQ) DeclareExchange(name, kind string) error {
	return r.channel.ExchangeDeclare(
		name,  // name
		kind,  // type (direct, fanout, topic, headers)
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		nil,   // arguments
	)
}

// DeclareQueue
func (r *RabbitMQ) DeclareQueue(name string) (amqp.Queue, error) {
	return r.channel.QueueDeclare(
		name,  // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
}

// BindQueue
func (r *RabbitMQ) BindQueue(queueName, routingKey, exchangeName string) error {
	return r.channel.QueueBind(
		queueName,    // queue name
		routingKey,   // routing key
		exchangeName, // exchange
		false,
		nil,
	)
}

// SetupInfrastructure
func (r *RabbitMQ) SetupInfrastructure() error {
	// Declare Exchange untuk IKAS events
	if err := r.DeclareExchange("ikas.events", "topic"); err != nil {
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
		if _, err := r.DeclareQueue(queueName); err != nil {
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
		if err := r.BindQueue(queueName, routingKey, "ikas.events"); err != nil {
			return fmt.Errorf("failed to bind queue %s: %w", queueName, err)
		}
	}

	log.Println("RabbitMQ infrastructure setup completed")
	return nil
}

// Reconnect ke RabbitMQ jika koneksi terputus
func (r *RabbitMQ) Reconnect() error {
	r.Close()
	time.Sleep(5 * time.Second)
	return r.connect()
}

// IsConnected
func (r *RabbitMQ) IsConnected() bool {
	return r.conn != nil && !r.conn.IsClosed()
}
