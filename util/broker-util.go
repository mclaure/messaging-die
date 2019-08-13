package util

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func getBrokerURL(url, username, password string) string {
	return fmt.Sprintf("amqp://%s:%s@%s", username, password, url)
}

// GetChannel ...
func GetChannel(url, username, password string) (*amqp.Connection, *amqp.Channel) {
	conn, err := amqp.Dial(getBrokerURL(url, username, password))
	failOnError(err, "Failed to establish connection to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to get an AMQP Channel for the connection")

	return conn, ch
}

// GetRPCQueue ...
func GetRPCQueue(queueName string, ch *amqp.Channel) *amqp.Queue {
	q, err := ch.QueueDeclare(
		queueName, // name string,
		false,     // durable bool,
		false,     // autoDelete bool,
		false,     // exclusive bool,
		false,     // noWait bool,
		nil)       // args amqp.Table)

	failOnError(err, fmt.Sprintf("Failed to declare queue: %s", queueName))

	// More details here:
	// https://www.rabbitmq.com/consumer-prefetch.html
	err = ch.Qos(
		1,     // prefetchCount int
		0,     // prefetchSize int
		false, // global bool
	)
	failOnError(err, "Failed to configure the Consumer Prefetch")

	return &q
}

// GetExclusiveQueue ...
func GetExclusiveQueue(ch *amqp.Channel) *amqp.Queue {
	q, err := ch.QueueDeclare(
		"",    // name string
		false, // durable bool
		false, // autoDelete bool
		true,  // exclusive bool
		false, // noWait bool
		nil)   // args amqp.Table

	failOnError(err, "Failed to declare an exclusive queue")

	return &q
}

// ConsumeFromRPCQueue ...
func ConsumeFromRPCQueue(queueName string, ch *amqp.Channel) <-chan amqp.Delivery {
	return consumeFromQueue(queueName, ch, false)
}

// ConsumeFromExclusiveQueue ...
func ConsumeFromExclusiveQueue(queueName string, ch *amqp.Channel) <-chan amqp.Delivery {
	return consumeFromQueue(queueName, ch, true)
}

// PublishToQueue ...
func PublishToQueue(queueName string, ch *amqp.Channel, msg amqp.Publishing) error {
	log.Printf("[*] Publishing to queue (%s)", queueName)
	err := ch.Publish(
		"",        // exchange string
		queueName, // key string
		false,     // mandatory bool
		false,     // immediate bool
		msg,       // msg amqp.Publishing
	)

	failOnError(err, fmt.Sprintf("Failed to publish a message on queue: %s", queueName))

	return err
}

func consumeFromQueue(queueName string, ch *amqp.Channel, autoAck bool) <-chan amqp.Delivery {
	log.Printf("[*] Consuming from queue (%s)", queueName)
	msgs, err := ch.Consume(
		queueName, // queue string
		"",        // consumer string
		autoAck,   // auto-ack bool
		false,     // exclusive bool
		false,     // no-local bool
		false,     // no-wait bool
		nil,       // args amqp.Table
	)

	failOnError(err, fmt.Sprintf("Failed to consume messages from queue: %s", queueName))

	return msgs
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)

		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
