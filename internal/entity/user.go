package entity

import "os"

type User struct {
	ID       int    `gorm:"primaryKey;not null;autoIncrement" json:"id"`
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (User) TableName() string {
	return os.Getenv("DB_PREFIX") + "user"
}
