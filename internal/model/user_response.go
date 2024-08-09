package model

import (
	"boiler-plate-clean/internal/entity"
	"time"
)

type UserResponse struct {
	ID        int             `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Email     string          `validate:"required,gt=2" json:"email,omitempty"`
	Password  string          `json:"password,omitempty"`
	CreatedAt *time.Time      `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
	Wallet    []entity.Wallet `json:"wallet"`
}

type UserRecapWallet struct {
	ID        int                `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Email     string             `validate:"required,gt=2" json:"email,omitempty"`
	Password  string             `json:"password,omitempty"`
	CreatedAt *time.Time         `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time         `gorm:"autoUpdateTime" json:"updated_at"`
	Wallet    []UserWalletDetail `json:"wallet"`
}

type UserWalletDetail struct {
	ID                   int                    `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Name                 string                 `json:"name"`
	TotalIncome          float64                `json:"total_income"`
	TotalExpense         float64                `json:"total_expense"`
	UserWalletTypeDetail []UserWalletTypeDetail `json:"user_wallet_type_detail"`
}

type UserWalletTypeDetail struct {
	Total       float64              `json:"total"`
	Type        string               `json:"type"`
	Transaction []entity.Transaction `json:"transaction"`
}
