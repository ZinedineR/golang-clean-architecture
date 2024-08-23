package service

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/model"
	"context"
	"github.com/RumbiaID/pkg-library/app/pkg/exception"
)

type UserService interface {
	// CRUD operations for User
	CreateUser(
		ctx context.Context, model *entity.User,
	) *exception.Exception
	FindById(ctx context.Context, id string) (*entity.User, *exception.Exception)
	Find(ctx context.Context, req *model.ListReq) (*[]entity.User, *exception.Exception)
}

type ListExampleResp struct {
	Pagination *model.Pagination `json:"pagination"`
	Data       []*entity.User    `json:"data"`
}
