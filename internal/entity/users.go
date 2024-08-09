package entity

import (
	"os"
	"time"
)

const (
	UsersTableName = "users"
)

type Users struct {
	ID        int        `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Email     string     `validate:"required,gt=2" json:"email,omitempty"`
	Password  string     `json:"password,omitempty"`
	CreatedAt *time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (model *Users) TableName() string {
	return os.Getenv("DB_PREFIX") + UsersTableName
}
