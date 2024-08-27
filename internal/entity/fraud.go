package entity

import (
	"os"
	"time"
)

const (
	FraudTableName = "fraud"
)

type Fraud struct {
	ID        int        `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	WalletId  int        `json:"wallet_id"`
	Wallet    *Wallet    `gorm:"foreignKey:WalletId;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"wallet,omitempty"`
	FraudTime *time.Time `gorm:"autoCreateTime" json:"fraud_time"`
}

func (model *Fraud) TableName() string {
	return os.Getenv("DB_PREFIX") + FraudTableName
}
