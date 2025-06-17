package main

import (
	"log"
	"os"

	"github.com/rabbitmq/amqp091-go"
)

func main() {
	log.Print(os.Getenv("RABBITMQ_URL"))
	// conn, err := amqp091.Dial(os.Getenv("RABBITMQ_URL")) // e.g., amqp://guest:guest@localhost:5672/
	conn, err := amqp091.Dial("amqp://guest:guest@localhost:5672/") // e.g., amqp://guest:guest@localhost:5672/
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to open a channel: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs",   // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("❌ Failed to declare exchange: %v", err)
	}

	q, err := ch.QueueDeclare(
		"",    // let RabbitMQ assign a random name
		false, // durable
		false, // delete when unused
		true,  // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("❌ Failed to declare queue: %v", err)
	}

	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		"logs", // exchange
		false,
		nil,
	)
	if err != nil {
		log.Fatalf("❌ Failed to bind queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		true,   // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		log.Fatalf("❌ Failed to register a consumer: %v", err)
	}

	log.Println("📥 Listening for logs. Press Ctrl+C to stop.")
	for msg := range msgs {
		log.Printf("📝 Log Received: %s", msg.Body)
	}
}
