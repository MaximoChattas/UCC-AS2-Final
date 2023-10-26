package queue

import (
	"Search/dto"
	"encoding/json"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
)

func Consume() (dto.QueueMessagesDto, error) {

	conn, err := amqp.Dial("amqp://user:password@rabbitmq:5672/")

	if err != nil {
		log.Debug("Failed to connect to RabbitMQ")
		return dto.QueueMessagesDto{}, err
	}

	defer conn.Close()

	channel, err := conn.Channel()

	if err != nil {
		log.Debug("Failed to open channel")
		return dto.QueueMessagesDto{}, err
	}

	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"hotel",
		false,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		log.Debug("Fail to declare a queue")
		return dto.QueueMessagesDto{}, err
	}

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
		log.Debug("Failed to publish consumer", err)
		return dto.QueueMessagesDto{}, err
	}

	var messages dto.QueueMessagesDto

	for msg := range msgs {

		var jsonMessage dto.QueueMessageDto

		err = json.Unmarshal(msg.Body, &jsonMessage)

		if err != nil {
			return dto.QueueMessagesDto{}, err
		}

		messages = append(messages, jsonMessage)
	}

	return messages, nil
}
