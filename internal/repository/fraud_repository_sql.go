package repository

import (
	"boiler-plate-clean/internal/entity"
)

type FraudRepo struct {
	Repository[entity.Fraud]
}

func NewFraudRepository() FraudRepository {
	return &FraudRepo{}
}
