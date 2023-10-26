package queue

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	"time"
)

func Publish(body []byte) error {

	conn, err := amqp.Dial("amqp://user:password@rabbitmq:5672/")

	if err != nil {
		log.Debug("Failed to connect to RabbitMQ")
		return err
	}

	defer conn.Close()

	channel, err := conn.Channel()

	if err != nil {
		log.Debug("Failed to open channel")
		return err
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
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = channel.PublishWithContext(
		ctx,
		"",
		queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})

	if err != nil {
		log.Debug("Error while publishing message", err)
		return err
	}

	return nil
}
