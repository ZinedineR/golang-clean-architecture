package service

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/pkg/exception"
	"context"
)

type TransactionService interface {
	// CRUD operations for Example
	Create(
		ctx context.Context, model *entity.Transaction,
	) *exception.Exception
	Update(
		ctx context.Context, id int, model *entity.Transaction,
	) *exception.Exception
	Detail(ctx context.Context, id int) (*entity.Transaction, *exception.Exception)
	Credit(
		ctx context.Context, walletid, categoryid int, amount float64,
	) *exception.Exception
	Transfer(
		ctx context.Context, senderId, receiverId int, amount float64,
	) *exception.Exception
	Delete(ctx context.Context, id int) *exception.Exception
}
