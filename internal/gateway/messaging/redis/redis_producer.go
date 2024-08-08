package redis

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"log/slog"
)

type ProducerRedis[T any] struct {
	Topic         string
	RedisProducer *redis.Client
}

func (p *ProducerRedis[T]) GetTopic() string {
	return p.Topic
}

func (p *ProducerRedis[T]) Send(ctx context.Context, data ...T) error {
	for _, d := range data {
		payload, err := json.Marshal(d)
		if err != nil {
			slog.Error("error when marshal payload", slog.String("error", err.Error()))
			return err
		}
		if err := p.RedisProducer.Publish(ctx, p.GetTopic(), payload); err != nil {
			slog.Error("failed to write messages", slog.String("error", err.String()))
			return err.Err()
		}
	}

	slog.Info("Success publish transaction", "key", p.Topic)

	return nil
}
