package queue

import (
	"Search/dto"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

var queue amqp.Queue
var channel *amqp.Channel

func InitQueue() {
	conn, err := amqp.Dial("amqp://user:password@localhost:5672/")
	if err != nil {
		log.Info("Failed to connect to RabbitMQ")
		log.Fatal(err)
	} else {
		log.Info("RabbitMQ connection established")
	}

	channel, err = conn.Channel()
	if err != nil {
		log.Info("Failed to open channel")
		log.Fatal(err)
	} else {
		log.Info("Channel opened")
	}

	queue, err = channel.QueueDeclare(
		"hotel",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Info("Failed to declare a queue")
		log.Fatal(err)
	} else {
		log.Info("Queue declared")
	}
}

func Consume() {

	msgs, err := channel.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		true,
		nil,
	)
	if err != nil {
		log.Error("Failed to publish consumer", err)
	}

	for msg := range msgs {

		var jsonMessage dto.QueueMessageDto

		err = json.Unmarshal(msg.Body, &jsonMessage)

		if err != nil {
			log.Error("Error:", err)
		}
	}
}
