package messaging

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log/slog"
)

type ConsumerRedisHandler func(ctx context.Context, message *redis.Message) error

func ConsumeRedisPublisher(
	ctx context.Context, redisService *redis.Client, key string,
	handler ConsumerRedisHandler,
) {
	reader := redisService.Subscribe(ctx, key)

	run := true

	for run {
		select {
		case <-ctx.Done():
			run = false
		default:
			message, err := reader.ReceiveMessage(ctx)
			if err == nil {
				err := handler(ctx, message)
				if err != nil {
					slog.Error("Failed to process message: ", slog.Any("error", err))
				}
			} else {
				slog.Warn("Consumer error: ", slog.Any("error", err))
			}
		}
	}

	slog.Info("Closing consumer for key : ", slog.String("key", key))
	err := reader.Close()
	if err != nil {
		panic(err)
	}
}
