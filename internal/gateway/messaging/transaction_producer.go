package messaging

import (
	"boiler-plate-clean/internal/entity"
	"context"
)

type TransactionProducer interface {
	GetTopic() string
	Send(ctx context.Context, order ...*entity.Transaction) error
}
