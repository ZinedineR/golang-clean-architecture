package repository

import (
	"boiler-plate-clean/internal/entity"
)

type TransactionRepo struct {
	Repository[entity.Transaction]
}

func NewTransactionRepository() TransactionRepository {
	return &TransactionRepo{}
}
