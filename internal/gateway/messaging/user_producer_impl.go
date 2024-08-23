package messaging

import (
	"boiler-plate-clean/internal/model"
	kafkaserver "github.com/RumbiaID/pkg-library/app/pkg/broker/kafkaservice"
)

type UserProducerImpl struct {
	ProducerKafka[*model.UserMessage]
}

func NewUserWriteProducerImpl(producer *kafkaserver.KafkaService, topic string) UserProducer {
	return &UserProducerImpl{
		ProducerKafka: ProducerKafka[*model.UserMessage]{
			Topic:         topic,
			KafkaProducer: producer,
		},
	}
}
