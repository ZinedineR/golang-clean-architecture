package service

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/model"
	"boiler-plate-clean/internal/repository"
	"boiler-plate-clean/pkg/exception"
	"context"
	"github.com/RumbiaID/pkg-library/app/pkg/xvalidator"
	"gorm.io/gorm"
	"time"
)

type UserServiceImpl struct {
	db              *gorm.DB
	userRepository  repository.UserRepository
	walletRepo      repository.WalletRepository
	transactionRepo repository.TransactionRepository
	validate        *xvalidator.Validator
}

func NewUserService(
	db *gorm.DB, repo repository.UserRepository,
	walletRepo repository.WalletRepository,
	transactionRepo repository.TransactionRepository,
	validate *xvalidator.Validator,
) UserService {
	return &UserServiceImpl{
		db:              db,
		userRepository:  repo,
		walletRepo:      walletRepo,
		transactionRepo: transactionRepo,
		validate:        validate,
	}
}

// CreateExample creates a new campaign
func (s *UserServiceImpl) Create(
	ctx context.Context, model *entity.Users,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()

	txRead := s.db

	result, err := s.userRepository.FindByName(ctx, txRead, "email", model.Email)
	if err != nil {
		return exception.Internal("err", err)
	}

	if result != nil {
		return exception.AlreadyExists("user already exists")
	}

	if err := s.userRepository.CreateTx(ctx, tx, model); err != nil {
		return exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *UserServiceImpl) Update(
	ctx context.Context, id int, model *entity.Users,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()
	txRead := s.db
	model.ID = id
	result, err := s.userRepository.FindByName(ctx, txRead, "email", model.Email)
	if err != nil {
		return exception.Internal("err", err)
	}

	if result != nil && result.ID != id {
		return exception.AlreadyExists("user already exists")
	}

	if err := s.userRepository.UpdateTx(ctx, tx, model); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *UserServiceImpl) Detail(ctx context.Context, id int) (*model.UserResponse, *exception.Exception) {
	result, err := s.userRepository.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if result == nil {
		return nil, exception.NotFound("user not found")
	}
	userResponse := &model.UserResponse{
		ID:        result.ID,
		Email:     result.Email,
		Password:  result.Password,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}
	wallets, err := s.walletRepo.FindByForeignKey(ctx, s.db, "user_id", result.ID)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if len(*wallets) > 0 {
		for _, wallet := range *wallets {
			userResponse.Wallet = append(userResponse.Wallet, wallet)
		}
	}
	return userResponse, nil
}

func (s *UserServiceImpl) Cashflow(ctx context.Context, id int, from time.Time, to time.Time) (
	*model.UserRecapWallet, *exception.Exception,
) {
	result, err := s.userRepository.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if result == nil {
		return nil, exception.NotFound("user not found")
	}
	userResponse := &model.UserRecapWallet{
		ID:        result.ID,
		Email:     result.Email,
		Password:  result.Password,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
	}

	wallets, err := s.walletRepo.FindByForeignKey(ctx, s.db, "user_id", result.ID)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	if len(*wallets) > 0 {
		for _, wallet := range *wallets {
			walletDetail := model.UserWalletDetail{
				ID:   wallet.ID,
				Name: wallet.Name,
			}
			incomeRecapWallet := model.UserWalletTypeDetail{
				Type: "income",
			}
			expenseRecapWallet := model.UserWalletTypeDetail{
				Type: "expense",
			}
			incomeTransactions, err := s.transactionRepo.FindByForeignKeyAndBetweenTimeWithFilter(ctx, s.db, "wallet_id", wallet.ID,
				"transaction_time", from, to,
				"type", "income")
			if err != nil {
				return nil, exception.Internal("err", err)
			}
			if len(*incomeTransactions) > 0 {
				for _, transaction := range *incomeTransactions {
					incomeRecapWallet.Total += transaction.Amount
					incomeRecapWallet.Transaction = append(incomeRecapWallet.Transaction, transaction)
				}
			}
			expenseTransactions, err := s.transactionRepo.FindByForeignKeyAndBetweenTimeWithFilter(ctx, s.db, "wallet_id", wallet.ID,
				"transaction_time", from, to,
				"type", "expense")
			if err != nil {
				return nil, exception.Internal("err", err)
			}
			if len(*expenseTransactions) > 0 {
				for _, transaction := range *expenseTransactions {
					expenseRecapWallet.Total += transaction.Amount
					expenseRecapWallet.Transaction = append(expenseRecapWallet.Transaction, transaction)
				}
			}
			walletDetail.TotalIncome = incomeRecapWallet.Total
			walletDetail.TotalExpense = expenseRecapWallet.Total
			walletDetail.UserWalletTypeDetail = append(walletDetail.UserWalletTypeDetail, incomeRecapWallet)
			walletDetail.UserWalletTypeDetail = append(walletDetail.UserWalletTypeDetail, expenseRecapWallet)
			userResponse.Wallet = append(userResponse.Wallet, walletDetail)
		}
	}
	return userResponse, nil
}

func (s *UserServiceImpl) Delete(ctx context.Context, id int) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()

	if err := s.userRepository.DeleteByIDTx(ctx, tx, id); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}
