package service

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/model"
	"boiler-plate-clean/pkg/exception"
	"context"
	"time"
)

type WalletService interface {
	// CRUD operations for Example
	Create(
		ctx context.Context, model *entity.Wallet,
	) *exception.Exception
	Update(
		ctx context.Context, id int, model *entity.Wallet,
	) *exception.Exception
	DetailWalletTransaction(ctx context.Context, id int, from, to time.Time) (
		*model.WalletResponse, *exception.Exception,
	)
	Detail(ctx context.Context, id int) (*entity.Wallet, *exception.Exception)
	Last10(ctx context.Context, id int) (*model.WalletResponse, *exception.Exception)
	RecapCategory(ctx context.Context, id, category int, from, to time.Time) (
		*model.WalletRecapCategoryResponse, *exception.Exception,
	)
	Delete(ctx context.Context, id int) *exception.Exception
}
