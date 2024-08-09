package entity

import (
	"os"
	"time"
)

const (
	TransactionTableName = "transaction"
)

type Transaction struct {
	ID                  int                  `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Type                string               `json:"type" validate:"eq=income|eq=expense|eq=transfer"`
	Amount              float64              `json:"amount"`
	Description         string               `json:"description"`
	WalletId            int                  `json:"wallet_id"`
	Wallet              *Wallet              `gorm:"foreignKey:WalletId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"wallet,omitempty"`
	CategoryId          int                  `json:"category_id"`
	CategoryTransaction *CategoryTransaction `gorm:"foreignKey:CategoryId;constraint:OnDelete:SET NULL;" json:"category_transaction,omitempty"`
	TransactionTime     *time.Time           `gorm:"autoCreateTime" json:"transaction_time"`
}

func (model *Transaction) TableName() string {
	return os.Getenv("DB_PREFIX") + TransactionTableName
}
