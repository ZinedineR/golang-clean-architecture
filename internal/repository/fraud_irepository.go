package repository

import (
	"boiler-plate-clean/internal/entity"
	"context"
	"gorm.io/gorm"
)

type FraudRepository interface {
	// Example operations
	CreateTx(ctx context.Context, tx *gorm.DB, data *entity.Fraud) error
}
