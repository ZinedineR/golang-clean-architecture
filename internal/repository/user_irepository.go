package repository

import (
	"boiler-plate-clean/internal/entity"
	"context"
	"gorm.io/gorm"
)

type UserRepository interface {
	// Example operations
	CreateTx(ctx context.Context, tx *gorm.DB, data *entity.Users) error
	UpdateTx(ctx context.Context, tx *gorm.DB, data *entity.Users) error
	FindByName(ctx context.Context, tx *gorm.DB, column, value string) (
		*entity.Users, error,
	)
	FindByID(ctx context.Context, tx *gorm.DB, id int) (*entity.Users, error)
	DeleteByIDTx(ctx context.Context, tx *gorm.DB, id int) error
}
