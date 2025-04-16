package repository

import (
	"boiler-plate-clean/internal/model"
	"boiler-plate-clean/pkg/pagination"
	"context"
	"errors"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"log/slog"
)

type BaseRepository[T any] interface {
	CreateTx(ctx context.Context, tx *gorm.DB, data *T) error
	UpdateTx(ctx context.Context, tx *gorm.DB, data *T) error
	UpdateTxWithAssociations(ctx context.Context, tx *gorm.DB, data *T) error
	DeleteByIDTx(ctx context.Context, tx *gorm.DB, id string) error
	Find(
		ctx context.Context, tx *gorm.DB, order model.OrderParam, filter model.FilterParams,
	) (*[]T, error)
	FindByPagination(
		ctx context.Context, tx *gorm.DB, page model.PaginationParam, order model.OrderParam,
		filter model.FilterParams,
	) (*model.PaginationData[T], error)
	FindByID(ctx context.Context, tx *gorm.DB, id string) (*T, error)
	FindByColumn(
		ctx context.Context, tx *gorm.DB, filter model.FilterParams, order model.OrderParam,
	) (*T, error)
}

type RelationField struct {
	Name string
	Func func(*gorm.DB) *gorm.DB
}

type BaseRepositoryImpl[T any] struct {
	relationFields []RelationField
}

func NewBaseRepositoryImpl[T any](
	relationFields []RelationField,
) BaseRepository[T] {
	return &BaseRepositoryImpl[T]{
		relationFields: relationFields,
	}
}

func (r *BaseRepositoryImpl[T]) CreateTx(ctx context.Context, tx *gorm.DB, data *T) error {
	if err := tx.WithContext(ctx).Omit(clause.Associations).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(data).Error; err != nil {
		slog.Error("failed to create", err)
		return err
	}
	return nil
}

func (r *BaseRepositoryImpl[T]) CreateTxAssociation(ctx context.Context, tx *gorm.DB, data *T) error {
	if err := tx.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			UpdateAll: true,
		}).
		Create(data).Error; err != nil {
		slog.Error("failed to create", err)
		return err
	}
	return nil
}

func (r *BaseRepositoryImpl[T]) UpdateTx(ctx context.Context, tx *gorm.DB, data *T) error {
	if err := tx.WithContext(ctx).Omit(clause.Associations).Model(data).Select("*").Updates(data).Error; err != nil {
		slog.Error("failed to update", err)
		return err
	}
	return nil
}

func (r *BaseRepositoryImpl[T]) UpdateTxWithAssociations(ctx context.Context, tx *gorm.DB, data *T) error {
	if err := tx.WithContext(ctx).Session(&gorm.Session{FullSaveAssociations: true}).Model(data).Select("*").Updates(data).Error; err != nil {
		slog.Error("failed to update", slog.Any("error", err))
		return err
	}
	return nil
}

func (r *BaseRepositoryImpl[T]) DeleteByIDTx(ctx context.Context, tx *gorm.DB, id string) error {
	if err := tx.WithContext(ctx).Unscoped().Where("id = ?", id).Delete(new(T)).Error; err != nil {
		slog.Error("failed to delete", err)
		return err
	}
	return nil
}

func (r *BaseRepositoryImpl[T]) FindByPagination(
	ctx context.Context, tx *gorm.DB, page model.PaginationParam, order model.OrderParam,
	filter model.FilterParams,
) (*model.PaginationData[T], error) {
	query := tx.WithContext(ctx).Omit(clause.Associations)
	query = pagination.Where(filter, query)
	query = pagination.Order(order, query)
	result, err := pagination.Paginate[T](page.Page, page.PageSize, query)
	if err != nil {
		return nil, err
	}
	return &model.PaginationData[T]{
		Page:             result.Page,
		PageSize:         result.PageSize,
		TotalPage:        result.TotalPage,
		TotalDataPerPage: result.TotalDataPerPage,
		TotalData:        result.TotalData,
		Data:             result.Data,
	}, nil
}

func (r *BaseRepositoryImpl[T]) Find(
	ctx context.Context, tx *gorm.DB, order model.OrderParam, filter model.FilterParams,
) (*[]T, error) {
	var data *[]T
	query := tx.WithContext(ctx).Omit(clause.Associations)
	query = pagination.Where(filter, query)
	query = pagination.Order(order, query)
	if err := query.Find(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("failed to find all", err)
		return nil, err
	}
	return data, nil
}

func (r *BaseRepositoryImpl[T]) FindByID(ctx context.Context, tx *gorm.DB, id string) (*T, error) {
	var data T
	if len(r.relationFields) > 0 {
		tx = tx.WithContext(ctx)
		for _, field := range r.relationFields {
			tx = tx.Preload(field.Name, field.Func)
		}
	} else {
		tx = tx.WithContext(ctx).Preload(clause.Associations)
	}
	if err := tx.Where("id = ?", id).First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("failed to find by id", err)
		return nil, err
	}
	return &data, nil
}

func (r *BaseRepositoryImpl[T]) FindByColumn(
	ctx context.Context, tx *gorm.DB, filter model.FilterParams, order model.OrderParam,
) (*T, error) {
	var data T
	query := tx.WithContext(ctx).Omit(clause.Associations)
	query = pagination.Where(filter, query)
	query = pagination.Order(order, query)
	if err := query.First(&data).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		slog.Error("failed to find by column", err)
		return nil, err
	}
	return &data, nil
}
