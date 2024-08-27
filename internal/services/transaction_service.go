package service

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/gateway/messaging"
	"boiler-plate-clean/internal/repository"
	"context"
	//"boiler-plate-clean/pkg/exception"
	"boiler-plate-clean/pkg/exception"
	"github.com/RumbiaID/pkg-library/app/pkg/xvalidator"
	"gorm.io/gorm"
	"strconv"
)

type TransactionServiceImpl struct {
	db                    *gorm.DB
	transactionRepository repository.TransactionRepository
	categoryRepository    repository.CategoryTransactionRepository
	userRepository        repository.UserRepository
	walletRepository      repository.WalletRepository
	transactionProducer   messaging.TransactionProducer
	validate              *xvalidator.Validator
}

func NewTransactionService(
	db *gorm.DB, repo repository.TransactionRepository,
	categoryRepository repository.CategoryTransactionRepository,
	userRepository repository.UserRepository,
	walletRepository repository.WalletRepository,
	transactionProducer messaging.TransactionProducer,
	validate *xvalidator.Validator,
) TransactionService {
	return &TransactionServiceImpl{
		db:                    db,
		transactionRepository: repo,
		categoryRepository:    categoryRepository,
		userRepository:        userRepository,
		walletRepository:      walletRepository,
		transactionProducer:   transactionProducer,
		validate:              validate,
	}
}

// CreateExample creates a new campaign
func (s *TransactionServiceImpl) Create(
	ctx context.Context, model *entity.Transaction,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()

	category, err := s.categoryRepository.FindByID(ctx, s.db, model.CategoryId)
	if err != nil {
		return exception.Internal("error in finding category", err)
	}
	if category == nil {
		return exception.PermissionDenied("category does not exists")
	}
	wallet, err := s.walletRepository.FindByID(ctx, s.db, model.WalletId)
	if err != nil {
		return exception.Internal("failed getting wallet detail", err)
	}
	if wallet == nil {
		return exception.NotFound("wallet detail not found")
	}
	model.Wallet = wallet
	if err := s.transactionRepository.CreateTx(ctx, tx, model); err != nil {
		return exception.Internal("err", err)
	}

	if err := s.transactionProducer.Send(ctx, model); err != nil {
		return exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *TransactionServiceImpl) Update(
	ctx context.Context, id int, model *entity.Transaction,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()
	category, err := s.categoryRepository.FindByID(ctx, s.db, model.CategoryId)
	if err != nil {
		return exception.Internal("error in finding category", err)
	}
	if category == nil {
		return exception.PermissionDenied("category does not exists")
	}
	if err := s.transactionRepository.UpdateTx(ctx, tx, model); err != nil {
		return exception.Internal("err", err)
	}
	if err := s.transactionProducer.Send(ctx, model); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *TransactionServiceImpl) Detail(ctx context.Context, id int) (
	*entity.Transaction, *exception.Exception,
) {
	result, err := s.transactionRepository.FindByID(ctx, s.db, id)
	if err != nil {
		return nil, exception.Internal("err", err)
	}
	return result, nil
}

func (s *TransactionServiceImpl) Credit(
	ctx context.Context, walletid, categoryid int, amount float64,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()
	if amount < 1 {
		return exception.PermissionDenied("Input of amount must be greater than zero")
	}
	wallet, err := s.walletRepository.FindByID(ctx, s.db, walletid)
	if err != nil {
		return exception.Internal("failed getting wallet detail", err)
	}
	if wallet == nil {
		return exception.NotFound("wallet detail not found")
	}

	category, err := s.categoryRepository.FindByID(ctx, s.db, categoryid)
	if err != nil {
		return exception.Internal("failed getting category detail", err)
	}
	if category == nil {
		return exception.PermissionDenied("category does not exists")
	}

	userTransaction := &entity.Transaction{
		Type:        "income",
		Amount:      amount,
		Description: "Credit of " + strconv.FormatFloat(amount, 'f', -1, 64),
		WalletId:    wallet.ID,
		Wallet:      wallet,
		CategoryId:  category.ID,
	}
	if err := s.transactionRepository.CreateTx(ctx, tx, userTransaction); err != nil {
		return exception.Internal("failed creating transaction", err)
	}

	wallet.Increase(amount)
	if err := s.walletRepository.UpdateTx(ctx, tx, wallet); err != nil {
		return exception.Internal("failed updating wallet", err)
	}
	if err := s.transactionProducer.Send(ctx, userTransaction); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *TransactionServiceImpl) Transfer(
	ctx context.Context, senderId, receiverId int, amount float64,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()
	if amount < 1 {
		return exception.PermissionDenied("Input of amount must be greater than zero")
	}
	sender, err := s.walletRepository.FindByID(ctx, s.db, senderId)
	if err != nil {
		return exception.Internal("failed getting sender detail", err)
	}
	if sender == nil {
		return exception.NotFound("sender wallet detail not found")
	}
	receiver, err := s.walletRepository.FindByID(ctx, s.db, receiverId)
	if err != nil {
		return exception.Internal("failed getting receiver detail", err)
	}
	if receiver == nil {
		return exception.NotFound("receiver wallet detail not found")
	}
	if sender.Balance < amount {
		return exception.PermissionDenied(sender.Name + " does not have enough balance. Balance: " + strconv.FormatFloat(sender.Balance, 'f', -1, 64))
	}
	category, err := s.categoryRepository.FindByName(ctx, s.db, "name", "Transfer")
	if err != nil {
		return exception.Internal("failed getting category detail", err)
	}
	if category == nil {
		return exception.PermissionDenied("category does not exists")
	}
	senderTransaction := &entity.Transaction{
		Type:        "transfer",
		Amount:      amount,
		Description: "Transfer to: " + receiver.Name,
		WalletId:    sender.ID,
		Wallet:      sender,
		CategoryId:  category.ID,
	}
	if err := s.transactionRepository.CreateTx(ctx, tx, senderTransaction); err != nil {
		return exception.Internal("failed creating transaction", err)
	}
	sender.Decrease(amount)
	if err := s.walletRepository.UpdateTx(ctx, tx, sender); err != nil {
		return exception.Internal("failed updating wallet", err)
	}
	receiverTransaction := &entity.Transaction{
		Type:        "transfer",
		Amount:      amount,
		Description: "Transfer from: " + sender.Name,
		WalletId:    receiver.ID,
		Wallet:      receiver,
		CategoryId:  category.ID,
	}
	if err := s.transactionRepository.CreateTx(ctx, tx, receiverTransaction); err != nil {
		return exception.Internal("failed creating transaction", err)
	}
	receiver.Increase(amount)
	if err := s.walletRepository.UpdateTx(ctx, tx, receiver); err != nil {
		return exception.Internal("failed updating wallet", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	if err := s.transactionProducer.Send(ctx, senderTransaction); err != nil {
		return exception.Internal("err", err)
	}
	if err := s.transactionProducer.Send(ctx, receiverTransaction); err != nil {
		return exception.Internal("err", err)
	}
	return nil
}

func (s *TransactionServiceImpl) Delete(ctx context.Context, id int) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()

	if err := s.transactionRepository.DeleteByIDTx(ctx, tx, id); err != nil {
		return exception.Internal("err", err)
	}
	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}
