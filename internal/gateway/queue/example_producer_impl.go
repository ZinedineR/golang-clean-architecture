package queue

import (
	"boiler-plate-clean/internal/model"
	"github.com/rabbitmq/amqp091-go"
)

type ExampleProducerImpl struct {
	Producer[*model.ExampleMessage]
}

func NewExampleProducerImpl(producer *amqp091.Connection, queue string) ExampleProducer {
	return &ExampleProducerImpl{
		Producer: Producer[*model.ExampleMessage]{
			QueueName:     queue,
			QueueProducer: producer,
		},
	}
}
