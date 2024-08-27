package service

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/pkg/exception"
	"context"
)

type FraudService interface {
	// CRUD operations for Example
	Create(
		ctx context.Context, model *entity.Fraud,
	) *exception.Exception
}
