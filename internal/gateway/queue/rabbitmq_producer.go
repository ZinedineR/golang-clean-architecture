package queue

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"log/slog"
)

type Producer[T any] struct {
	QueueName     string
	QueueProducer *amqp091.Connection
}

func (p *Producer[T]) GetQueue() string {
	return p.QueueName
}

func (p *Producer[T]) Send(ctx context.Context, queueName string, data ...T) error {
	ch, err := p.QueueProducer.Channel()
	if err != nil {
		logrus.Print("Failed to create RabbitMQ channel")
	}
	defer ch.Close()

	queue, err := declareQueue(ch, queueName)
	if err != nil {
		return err
	}
	for _, d := range data {
		payload, err := json.Marshal(d)
		if err != nil {
			slog.Error("error when marshal payload", slog.String("error", err.Error()))
			return err
		}
		if err := ch.PublishWithContext(ctx,
			"",
			queue.Name,
			false,
			false,
			amqp091.Publishing{
				ContentType: "application/json",
				Body:        payload,
			},
		); err != nil {
			slog.Error("Error publish message to queue: ", queue.Name)
			return err
		}
	}

	slog.Info("Success publish queue, ", "key", p.QueueName)

	return nil
}

func declareQueue(ch *amqp091.Channel, queueName string) (*amqp091.Queue, error) {
	queue, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		slog.Error("Failed to declare %s queue", queueName)
		return nil, err
	}
	return &queue, nil
}
