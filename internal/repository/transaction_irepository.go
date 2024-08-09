package repository

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/model"
	"context"
	"gorm.io/gorm"
	"time"
)

type TransactionRepository interface {
	// Example operations
	CreateTx(ctx context.Context, tx *gorm.DB, data *entity.Transaction) error
	UpdateTx(ctx context.Context, tx *gorm.DB, data *entity.Transaction) error
	FindByName(ctx context.Context, tx *gorm.DB, column, value string) (
		*entity.Transaction, error,
	)
	Find(
		ctx context.Context, tx *gorm.DB, order model.OrderParam, filter model.FilterParams,
	) (*[]entity.Transaction, error)
	FindByPagination(
		ctx context.Context, tx *gorm.DB, page model.PaginationParam, order model.OrderParam,
		filter model.FilterParams,
	) (*model.PaginationData[entity.Transaction], error)
	FindByForeignKeyAndBetweenTime(
		ctx context.Context, tx *gorm.DB, column string, id int,
		columnDate string, from, to time.Time,
	) (*[]entity.Transaction, error)
	FindAssociationByForeignKeyAndBetweenTime(
		ctx context.Context, tx *gorm.DB, column string, id int,
		columnDate string, from, to time.Time,
	) (*[]entity.Transaction, error)
	FindByForeignKeyAndBetweenTimeWithFilter(
		ctx context.Context, tx *gorm.DB, column string, id int,
		columnDate string, from, to time.Time,
		columnFilter string, value string,
	) (*[]entity.Transaction, error)
	FindByID(ctx context.Context, tx *gorm.DB, id int) (*entity.Transaction, error)
	DeleteByIDTx(ctx context.Context, tx *gorm.DB, id int) error
}
