package repository

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/model"
	"context"
	"gorm.io/gorm"
)

type UserRepository interface {
	// User operations
	CreateTx(ctx context.Context, tx *gorm.DB, data *entity.User) error
	FindByName(ctx context.Context, tx *gorm.DB, column, value string) (
		*entity.User, error,
	)
	FindByID(ctx context.Context, tx *gorm.DB, id string) (*entity.User, error)
	Find(
		ctx context.Context, tx *gorm.DB, order model.OrderParam, filter model.FilterParams,
	) (*[]entity.User, error)
}
