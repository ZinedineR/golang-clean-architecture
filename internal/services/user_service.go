package service

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/model"
	"boiler-plate-clean/internal/repository"
	"context"
	"github.com/RumbiaID/pkg-library/app/pkg/exception"
	"github.com/RumbiaID/pkg-library/app/pkg/xvalidator"
	"gorm.io/gorm"
	"strconv"
)

type UserServiceImpl struct {
	db           *gorm.DB
	campaignRepo repository.UserRepository
	validate     *xvalidator.Validator
}

func NewUserService(
	db *gorm.DB, repo repository.UserRepository,
	validate *xvalidator.Validator,
) UserService {
	return &UserServiceImpl{
		db:           db,
		campaignRepo: repo,
		validate:     validate,
	}
}

// CreateExample creates a new campaign
func (s *UserServiceImpl) CreateUser(
	ctx context.Context, model *entity.User,
) *exception.Exception {
	tx := s.db.Begin()
	defer tx.Rollback()

	txRead := s.db
	if errs := s.validate.Struct(model); errs != nil {
		return exception.InvalidArgument(errs)
	}
	result, err := s.campaignRepo.FindByName(ctx, txRead, "name", model.Name)
	if err != nil {
		return exception.Internal("err", err)
	}

	if result != nil {
		return exception.AlreadyExists("example already exists")
	}

	if err := s.campaignRepo.CreateTx(ctx, tx, model); err != nil {
		return exception.Internal("err", err)
	}

	if err := tx.Commit().Error; err != nil {
		return exception.Internal("commit transaction", err)
	}
	return nil
}

func (s *UserServiceImpl) FindById(ctx context.Context, id string) (*entity.User, *exception.Exception) {

	txRead := s.db
	_, err := strconv.Atoi(id)
	if err != nil {
		return nil, exception.Internal("id is not string", err)
	}
	result, err := s.campaignRepo.FindByID(ctx, txRead, id)
	if err != nil {
		return nil, exception.Internal("err", err)
	}

	return result, nil
}

func (s *UserServiceImpl) Find(ctx context.Context, req *model.ListReq) (*[]entity.User, *exception.Exception) {

	txRead := s.db
	result, err := s.campaignRepo.Find(ctx, txRead, req.Order, req.Filter)
	if err != nil {
		return nil, exception.Internal("err", err)
	}

	return result, nil
}
