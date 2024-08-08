package messaging

import (
	"boiler-plate-clean/internal/gateway/messaging/kafka"
	redisserver "boiler-plate-clean/internal/gateway/messaging/redis"
	"boiler-plate-clean/internal/model"
	kafkaserver "github.com/RumbiaID/pkg-library/app/pkg/broker/kafkaservice"
	"github.com/go-redis/redis/v8"
)

type ExampleProducerImpl struct {
	kafka.ProducerKafka[*model.ExampleMessage]
}

func NewExampleKafkaProducerImpl(producer *kafkaserver.KafkaService, topic string) ExampleProducer {
	return &ExampleProducerImpl{
		ProducerKafka: kafka.ProducerKafka[*model.ExampleMessage]{
			Topic:         topic,
			KafkaProducer: producer,
		},
	}
}

type EmailRedisProducerImpl struct {
	redisserver.ProducerRedis[*model.ExampleMessage]
}

func NewExampleRedisProducerImpl(producer *redis.Client, topic string) ExampleProducer {
	return &EmailRedisProducerImpl{
		ProducerRedis: redisserver.ProducerRedis[*model.ExampleMessage]{
			Topic:         topic,
			RedisProducer: producer,
		},
	}
}
