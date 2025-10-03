package connections

import (
	"encoding/json"
	"log"

	"github.com/streadway/amqp"
	"github.com/substrate-cli/api-server/internal/utils"
)

var channel *amqp.Channel
var exchangeName = "dev.topic.spinrequest"

func InitRabbitMQ() {
	conn, err := amqp.Dial(utils.GetAMQPUrl())
	if err != nil {
		log.Fatalf("❌ Failed to connect to RabbitMQ: %v", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("❌ Failed to open channel: %v", err)
	}

	// Declare exchange to be safe
	err = ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true, false, false, false, nil,
	)
	if err != nil {
		log.Fatalf("❌ Failed to declare exchange: %v", err)
	}

	channel = ch
	log.Println("✅ RabbitMQ publisher initialized")
}

func PublishSpinRequest(payload interface{}, routingKey string) error {
	messageBytes, err := json.Marshal(payload)
	if err != nil {
		log.Println("Failed to parse json")
		return err
	}
	err = channel.Publish(
		exchangeName, // exchange
		routingKey,   // routing key (e.g., "spin.create")
		false,        // mandatory
		false,        // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(messageBytes),
		},
	)
	if err != nil {
		log.Printf("❌ Failed to publish message: %v", err)
	}
	log.Printf("Published: [%s] %s", routingKey, messageBytes)
	return nil
}
