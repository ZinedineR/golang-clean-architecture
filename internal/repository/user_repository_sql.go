package repository

import (
	"boiler-plate-clean/internal/entity"
)

type UserRepo struct {
	Repository[entity.Users]
}

func NewUserRepository() UserRepository {
	return &UserRepo{}
}
