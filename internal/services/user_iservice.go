package service

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/model"
	"boiler-plate-clean/pkg/exception"
	"context"
	"time"
)

type UserService interface {
	// CRUD operations for Example
	Create(
		ctx context.Context, model *entity.Users,
	) *exception.Exception
	Update(
		ctx context.Context, id int, model *entity.Users,
	) *exception.Exception
	Detail(ctx context.Context, id int) (*model.UserResponse, *exception.Exception)
	Delete(ctx context.Context, id int) *exception.Exception
	Cashflow(ctx context.Context, id int, from time.Time, to time.Time) (*model.UserRecapWallet, *exception.Exception)
}

type ListExampleResp struct {
	Pagination *model.Pagination `json:"pagination"`
	Data       []*entity.Users   `json:"data"`
}
