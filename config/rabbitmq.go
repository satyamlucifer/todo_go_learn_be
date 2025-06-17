package config

import (
	"log"
	"os"

	amqp "github.com/rabbitmq/amqp091-go"
)

var RMQConn *amqp.Connection
var RMQChannel *amqp.Channel

func InitRabbitMQ() {
	var err error
	RMQConn, err = amqp.Dial("", os.Getenv("RABBITMQ_URL")) // "amqp://guest:guest@localhost:5672/"
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}

	RMQChannel, err = RMQConn.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to open channel: %v", err)
	}

	err = RMQChannel.ExchangeDeclare("logs", "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("❌ Failed to declare exchange: %v", err)
	}

	log.Println("✅ Connected to RabbitMQ")
}

func PublishLog(msg string) {
	if RMQChannel == nil {
		log.Println("⚠️ RMQ not initialized")
		return
	}
	RMQChannel.Publish("logs", "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
}
