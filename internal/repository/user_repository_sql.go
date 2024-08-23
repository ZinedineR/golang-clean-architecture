package repository

import (
	"boiler-plate-clean/internal/entity"
)

type UserRepo struct {
	Repository[entity.User]
}

func NewUserRepository() UserRepository {
	return &UserRepo{}
}
