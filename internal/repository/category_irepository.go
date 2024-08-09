package repository

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/model"
	"context"
	"gorm.io/gorm"
)

type CategoryTransactionRepository interface {
	// Example operations
	CreateTx(ctx context.Context, tx *gorm.DB, data *entity.CategoryTransaction) error
	UpdateTx(ctx context.Context, tx *gorm.DB, data *entity.CategoryTransaction) error
	FindByName(ctx context.Context, tx *gorm.DB, column, value string) (
		*entity.CategoryTransaction, error,
	)
	Find(
		ctx context.Context, tx *gorm.DB, order model.OrderParam, filter model.FilterParams,
	) (*[]entity.CategoryTransaction, error)
	FindByID(ctx context.Context, tx *gorm.DB, id int) (*entity.CategoryTransaction, error)
	DeleteByIDTx(ctx context.Context, tx *gorm.DB, id int) error
}
