package queue

import (
	"context"
	"encoding/json"
	"github.com/rabbitmq/amqp091-go"
	"log/slog"
)

type ExampleConsumer struct {
}

func NewExampleConsumer() *ExampleConsumer {
	return &ExampleConsumer{}
}

func (c ExampleConsumer) Consume(ctx context.Context, msg *amqp091.Delivery) error {
	data := new(any)
	if err := json.Unmarshal(msg.Body, data); err != nil {
		slog.Error("error unmarshalling notification event", slog.String("error", err.Error()))
		return err
	}
	slog.Info("Received topic notification with event", slog.Any("sms", data))
	//Simulate handlebars
	return nil
}
