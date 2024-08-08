package queue

import (
	"boiler-plate-clean/internal/model"
	"context"
)

type ExampleProducer interface {
	GetQueue() string
	Send(ctx context.Context, queueName string, data ...*model.ExampleMessage) error
}
