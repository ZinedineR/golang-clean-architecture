package repository

import (
	"boiler-plate-clean/internal/entity"
)

type CategoryTransactionRepo struct {
	Repository[entity.CategoryTransaction]
}

func NewCategoryTransactionRepository() CategoryTransactionRepository {
	return &CategoryTransactionRepo{}
}
