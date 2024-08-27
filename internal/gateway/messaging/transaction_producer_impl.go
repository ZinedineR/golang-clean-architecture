package messaging

import (
	"boiler-plate-clean/internal/entity"
	kafkaserver "github.com/RumbiaID/pkg-library/app/pkg/broker/kafkaservice"
)

type TransactionProducerImpl struct {
	ProducerKafka[*entity.Transaction]
}

func NewTransactionProducerImpl(producer *kafkaserver.KafkaService, topic string) TransactionProducer {
	return &TransactionProducerImpl{
		ProducerKafka: ProducerKafka[*entity.Transaction]{
			Topic:         topic,
			KafkaProducer: producer,
		},
	}
}
