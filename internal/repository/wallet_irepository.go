package repository

import (
	"boiler-plate-clean/internal/entity"
	"context"
	"gorm.io/gorm"
)

type WalletRepository interface {
	// Example operations
	CreateTx(ctx context.Context, tx *gorm.DB, data *entity.Wallet) error
	UpdateTx(ctx context.Context, tx *gorm.DB, data *entity.Wallet) error
	FindByName(ctx context.Context, tx *gorm.DB, column, value string) (
		*entity.Wallet, error,
	)
	FindByForeignKey(ctx context.Context, tx *gorm.DB, column string, id int) (*[]entity.Wallet, error)
	FindByID(ctx context.Context, tx *gorm.DB, id int) (*entity.Wallet, error)
	DeleteByIDTx(ctx context.Context, tx *gorm.DB, id int) error
}
