package service

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/repository"
	"boiler-plate-clean/pkg/exception"
	"context"
	"github.com/RumbiaID/pkg-library/app/pkg/xvalidator"
	"gorm.io/gorm"
)

type FraudServiceImpl struct {
	db              *gorm.DB
	fraudRepository repository.FraudRepository
	validate        *xvalidator.Validator
}

func NewFraudService(
	db *gorm.DB, repo repository.FraudRepository,
	validate *xvalidator.Validator,
) FraudService {
	return &FraudServiceImpl{
		db:              db,
		fraudRepository: repo,
		validate:        validate,
	}
}

// CreateExample creates a new campaign
func (s *FraudServiceImpl) Create(
	ctx context.Context, model *entity.Fraud,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()

	if err := s.fraudRepository.CreateTx(ctx, tx, model); err != nil {
		return exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}
