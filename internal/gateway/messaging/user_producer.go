package messaging

import (
	"boiler-plate-clean/internal/model"
	"context"
)

type UserProducer interface {
	GetTopic() string
	Send(ctx context.Context, order ...*model.UserMessage) error
}
