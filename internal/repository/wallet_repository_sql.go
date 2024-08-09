package repository

import (
	"boiler-plate-clean/internal/entity"
)

type WalletRepo struct {
	Repository[entity.Wallet]
}

func NewWalletRepository() WalletRepository {
	return &WalletRepo{}
}
