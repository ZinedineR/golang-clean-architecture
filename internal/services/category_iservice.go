package service

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/model"
	"boiler-plate-clean/pkg/exception"
	"context"
)

type CategoryTransactionService interface {
	// CRUD operations for Example
	Create(
		ctx context.Context, model *entity.CategoryTransaction,
	) *exception.Exception
	Update(
		ctx context.Context, id int, model *entity.CategoryTransaction,
	) *exception.Exception
	Detail(ctx context.Context, id int) (*entity.CategoryTransaction, *exception.Exception)
	Find(ctx context.Context, order model.OrderParam, filter model.FilterParams) (
		*[]entity.CategoryTransaction, *exception.Exception,
	)
	Delete(ctx context.Context, id int) *exception.Exception
}
