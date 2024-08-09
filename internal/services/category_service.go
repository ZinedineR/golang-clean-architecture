package service

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/model"
	"boiler-plate-clean/internal/repository"
	"boiler-plate-clean/pkg/exception"
	"context"
	"github.com/RumbiaID/pkg-library/app/pkg/xvalidator"
	"gorm.io/gorm"
)

type CategoryTransactionServiceImpl struct {
	db                 *gorm.DB
	categoryRepository repository.CategoryTransactionRepository
	validate           *xvalidator.Validator
}

func NewCategoryTransactionService(
	db *gorm.DB, repo repository.CategoryTransactionRepository,
	validate *xvalidator.Validator,
) CategoryTransactionService {
	return &CategoryTransactionServiceImpl{
		db:                 db,
		categoryRepository: repo,
		validate:           validate,
	}
}

// CreateExample creates a new campaign
func (s *CategoryTransactionServiceImpl) Create(
	ctx context.Context, model *entity.CategoryTransaction,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()

	if err := s.categoryRepository.CreateTx(ctx, tx, model); err != nil {
		return exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *CategoryTransactionServiceImpl) Update(
	ctx context.Context, id int, model *entity.CategoryTransaction,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()
	model.ID = id
	if err := s.categoryRepository.UpdateTx(ctx, tx, model); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *CategoryTransactionServiceImpl) Detail(ctx context.Context, id int) (
	*entity.CategoryTransaction, *exception.Exception,
) {
	result, err := s.categoryRepository.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	return result, nil
}

func (s *CategoryTransactionServiceImpl) Find(ctx context.Context, order model.OrderParam, filter model.FilterParams) (
	*[]entity.CategoryTransaction, *exception.Exception,
) {
	result, err := s.categoryRepository.Find(ctx, s.db, order, filter)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	return result, nil
}

func (s *CategoryTransactionServiceImpl) Delete(ctx context.Context, id int) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()

	if err := s.categoryRepository.DeleteByIDTx(ctx, tx, id); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}
