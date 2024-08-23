package messaging

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/model"
	service "boiler-plate-clean/internal/services"
	"context"
	"encoding/json"
	"errors"
	"github.com/segmentio/kafka-go"
	"log/slog"
)

type UserConsumer struct {
	UserWrite service.UserService
	UserRead  service.UserService
}

func NewUserWriteConsumer(userWrite, UserRead service.UserService) *UserConsumer {
	return &UserConsumer{
		UserWrite: userWrite,
		UserRead:  UserRead,
	}
}

func (c UserConsumer) ConsumeKafka(ctx context.Context, message *kafka.Message) error {
	userEvent := new(model.UserMessage)
	if err := json.Unmarshal(message.Value, userEvent); err != nil {
		slog.Error("error unmarshalling example event", slog.String("error", err.Error()))
		return err
	}
	if err := c.UserWrite.CreateUser(ctx, &entity.User{
		Name:     userEvent.Name,
		Password: userEvent.Password,
	}); err != nil {
		slog.Error("error creating user", slog.Any("error", err))
		return errors.New("error creating user")
	}

	if err := c.UserRead.CreateUser(ctx, &entity.User{
		Name:     userEvent.Name,
		Password: userEvent.Password,
	}); err != nil {
		slog.Error("error creating user", slog.Any("error", err))
		return errors.New("error creating user")
	}

	slog.Info("Received topic example with event", slog.Any("example", userEvent))
	return nil
}
