package queue

import (
	"context"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type ConsumerHandler func(ctx context.Context, message *amqp091.Delivery) error

func ConsumePublisher(
	ctx context.Context, rabbitConnection *amqp091.Connection, queueName string,
	handler ConsumerHandler,
) {
	ch, err := rabbitConnection.Channel()
	if err != nil {
		panic(err)
	}
	queue, err := declareQueue(ch, queueName)
	if err != nil {
		panic(err)
	}
	err = ch.Qos(8, 0, false)
	if err != nil {
		panic(err)
	}
	msgs, err := ch.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		slog.Error("Failed to register a consumer: ", slog.Any("error", err))
		panic(err)
	}

	run := true
	for run {
		select {
		case <-ctx.Done():
			run = false
		case msg := <-msgs:
			err := handler(ctx, &msg)
			if err != nil {
				slog.Error("Failed to process message: ", slog.Any("error", err))
				err := msg.Nack(false, false)
				if err != nil {
					return
				}
			} else {
				err := msg.Ack(false)
				if err != nil {
					return
				}
			}
		}
	}
	slog.Info("Closing consumer for queue : ", slog.Any("queue", queue))
	err = ch.Close()
	if err != nil {
		panic(err)
	}
}

func declareQueue(ch *amqp091.Channel, queueName string) (*amqp091.Queue, error) {
	queue, err := ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		slog.Error("Failed to declare queue: " + queueName)
		return nil, err
	}
	return &queue, nil
}
