package model

import (
	"boiler-plate-clean/internal/entity"
	"time"
)

type WalletResponse struct {
	ID              int                  `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	UserId          int                  `bson:"user_id" json:"user_id" validate:"required,uuid" gorm:"type:uuid"`
	Balance         float64              `gorm:"default:0" json:"balance"`
	LastTransaction *time.Time           `gorm:"autoUpdateTime" json:"last_transaction"`
	Transaction     []entity.Transaction `json:"transaction"`
}

type WalletRecapCategoryResponse struct {
	ID          int                         `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	UserId      int                         `bson:"user_id" json:"user_id" validate:"required,uuid" gorm:"type:uuid"`
	Total       float64                     `gorm:"default:0" json:"total"`
	CategoryId  int                         `bson:"category_id" json:"category_id"`
	Category    *entity.CategoryTransaction `json:"category"`
	Transaction []entity.Transaction        `json:"transaction"`
}
