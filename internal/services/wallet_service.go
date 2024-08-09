package service

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/model"
	"boiler-plate-clean/internal/repository"
	"boiler-plate-clean/pkg/exception"
	"context"
	"github.com/RumbiaID/pkg-library/app/pkg/xvalidator"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type WalletServiceImpl struct {
	db                 *gorm.DB
	userRepository     repository.UserRepository
	walletRepository   repository.WalletRepository
	transactionRepo    repository.TransactionRepository
	categoryRepository repository.CategoryTransactionRepository
	validate           *xvalidator.Validator
}

func NewWalletService(
	db *gorm.DB, repo repository.WalletRepository,
	userRepository repository.UserRepository,
	transactionRepository repository.TransactionRepository,
	categoryRepository repository.CategoryTransactionRepository,
	validate *xvalidator.Validator,
) WalletService {
	return &WalletServiceImpl{
		db:                 db,
		walletRepository:   repo,
		userRepository:     userRepository,
		transactionRepo:    transactionRepository,
		categoryRepository: categoryRepository,
		validate:           validate,
	}
}

// CreateExample creates a new campaign
func (s *WalletServiceImpl) Create(
	ctx context.Context, model *entity.Wallet,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()

	userCheck, err := s.userRepository.FindByID(ctx, s.db, model.UserId)
	if err != nil {
		return exception.Internal("error finding user", err)
	}
	if userCheck == nil {
		return exception.PermissionDenied("user does not exists")
	}

	duplicateCheck, err := s.walletRepository.FindByName(ctx, s.db, "name", model.Name)
	if err != nil {
		return exception.Internal("error finding wallet", err)
	}
	if duplicateCheck != nil && duplicateCheck.User.ID == userCheck.ID {
		return exception.PermissionDenied("wallet already exists")
	}
	if err := s.walletRepository.CreateTx(ctx, tx, model); err != nil {
		return exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *WalletServiceImpl) Update(
	ctx context.Context, id int, model *entity.Wallet,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()
	model.ID = id
	userCheck, err := s.userRepository.FindByID(ctx, s.db, model.UserId)
	if err != nil {
		return exception.Internal("error finding user", err)
	}
	if userCheck == nil {
		return exception.PermissionDenied("user does not exists")
	}

	duplicateCheck, err := s.walletRepository.FindByName(ctx, s.db, "name", model.Name)
	if err != nil {
		return exception.Internal("error finding wallet", err)
	}
	if duplicateCheck != nil && duplicateCheck.User.ID == userCheck.ID && duplicateCheck.ID != model.ID {
		return exception.PermissionDenied("wallet already exists")
	}
	if err := s.walletRepository.UpdateTx(ctx, tx, model); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *WalletServiceImpl) Detail(ctx context.Context, id int, from, to time.Time) (
	*model.WalletResponse, *exception.Exception,
) {
	result, err := s.walletRepository.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if result == nil {
		return nil, exception.PermissionDenied("wallet not found")
	}
	walletResponse := &model.WalletResponse{
		ID:              id,
		UserId:          result.UserId,
		Balance:         result.Balance,
		LastTransaction: result.LastTransaction,
	}
	transactions, err := s.transactionRepo.FindByForeignKeyAndBetweenTime(ctx, s.db, "wallet_id", result.ID,
		"transaction_time", from, to)
	if err != nil {
		return nil, exception.Internal("err", err)
	}

	if len(*transactions) > 0 {
		for _, transaction := range *transactions {
			walletResponse.Transaction = append(walletResponse.Transaction, transaction)
		}
	}

	return walletResponse, nil
}

func (s *WalletServiceImpl) Last10(ctx context.Context, id int) (
	*model.WalletResponse, *exception.Exception,
) {
	result, err := s.walletRepository.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if result == nil {
		return nil, exception.PermissionDenied("wallet not found")
	}
	walletResponse := &model.WalletResponse{
		ID:              id,
		UserId:          result.UserId,
		Balance:         result.Balance,
		LastTransaction: result.LastTransaction,
	}
	pageParam := model.PaginationParam{
		Page:     1,
		PageSize: 10,
	}
	orderParam := model.OrderParam{
		Order:   "desc",
		OrderBy: "transaction_time",
	}
	filterParam := model.FilterParams{
		{
			Field:    "wallet_id",
			Value:    strconv.Itoa(result.ID),
			Operator: "=",
		},
	}
	transactions, err := s.transactionRepo.FindByPagination(ctx, s.db, pageParam, orderParam, filterParam)
	if err != nil {
		return nil, exception.Internal("err", err)
	}

	if len(transactions.Data) > 0 {
		for _, transaction := range transactions.Data {
			walletResponse.Transaction = append(walletResponse.Transaction, *transaction)
		}
	}

	return walletResponse, nil
}

func (s *WalletServiceImpl) RecapCategory(ctx context.Context, id, category int, from, to time.Time) (
	*model.WalletRecapCategoryResponse, *exception.Exception,
) {
	result, err := s.walletRepository.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if result == nil {
		return nil, exception.PermissionDenied("wallet not found")
	}

	categoryCheck, err := s.categoryRepository.FindByID(ctx, s.db, category)
	if err != nil {
		return nil, exception.Internal("err", err)
	}

	if categoryCheck == nil {
		return nil, exception.PermissionDenied("wallet not found")
	}

	walletResponse := &model.WalletRecapCategoryResponse{
		ID:         id,
		UserId:     result.UserId,
		CategoryId: categoryCheck.ID,
		Category:   categoryCheck,
	}
	transactions, err := s.transactionRepo.FindAssociationByForeignKeyAndBetweenTime(ctx, s.db, "wallet_id", result.ID,
		"transaction_time", from, to)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if len(*transactions) > 0 {
		for _, transaction := range *transactions {
			if transaction.CategoryTransaction.Name == categoryCheck.Name && transaction.Type == "expense" {
				transaction.Wallet = nil
				transaction.CategoryTransaction = nil
				walletResponse.Total += transaction.Amount
				walletResponse.Transaction = append(walletResponse.Transaction, transaction)
			}
		}
	}

	return walletResponse, nil
}

func (s *WalletServiceImpl) Delete(ctx context.Context, id int) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()

	if err := s.walletRepository.DeleteByIDTx(ctx, tx, id); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}
