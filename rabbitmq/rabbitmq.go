package rabbitmq

import (
	"github.com/streadway/amqp"
	"os"
)

func NewConnection() (*amqp.Connection, error) {
	url := os.Getenv("RABBITMQ_URL")
	return amqp.Dial(url)
}

func PublishMessage(ch *amqp.Channel, queueName string, body []byte) error {
	_, err := ch.QueueDeclare(
		queueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
	return err
}
