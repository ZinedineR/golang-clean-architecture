package messaging

import (
	"boiler-plate-clean/internal/gateway/queue"
	"boiler-plate-clean/internal/model"
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

type ExampleConsumer struct {
	ExampleQueue queue.ExampleProducer
}

func NewExampleConsumer(queue queue.ExampleProducer) *ExampleConsumer {
	return &ExampleConsumer{
		ExampleQueue: queue,
	}
}

func (c ExampleConsumer) ConsumeKafka(ctx context.Context, message *kafka.Message) error {
	exampleEvent := new(model.ExampleMessage)
	if err := json.Unmarshal(message.Value, exampleEvent); err != nil {
		slog.Error("error unmarshalling example event", slog.String("error", err.Error()))
		return err
	}
	slog.Info("Received topic example with event", slog.Any("example", exampleEvent))
	if err := c.ExampleQueue.Send(ctx, "ExampleQueue", exampleEvent); err != nil {
		slog.Error("error sending example", slog.String("error", err.Error()))
	}
	return nil
}

func (c ExampleConsumer) ConsumeRedis(ctx context.Context, message *redis.Message) error {
	exampleEvent := new(model.ExampleMessage)
	if err := json.Unmarshal([]byte(message.Payload), exampleEvent); err != nil {
		slog.Error("error unmarshalling example event", slog.String("error", err.Error()))
		return err
	}

	slog.Info("Received topic example with event", slog.Any("example", exampleEvent))
	if err := c.ExampleQueue.Send(ctx, "ExampleQueue", exampleEvent); err != nil {
		slog.Error("error sending example", slog.String("error", err.Error()))
	}
	return nil
}
