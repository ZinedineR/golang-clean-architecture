package entity

import (
	"os"
	"time"
)

const (
	CategoryTransactionTableName = "category_transaction"
)

type CategoryTransaction struct {
	ID        int        `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Name      string     `validate:"required,gt=2" json:"name,omitempty"`
	CreatedAt *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (model *CategoryTransaction) TableName() string {
	return os.Getenv("DB_PREFIX") + CategoryTransactionTableName
}
