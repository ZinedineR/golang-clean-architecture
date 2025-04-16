package repository

import (
	"boiler-plate-clean/internal/entity"
	"gorm.io/gorm"
)

type ExampleSQLRepo struct {
	BaseRepository[entity.Example]
}

func NewExampleSQLRepository() ExampleRepository {
	relationField := []RelationField{
		{
			Name: "Example",
			Func: func(db *gorm.DB) *gorm.DB {
				return db
			},
		},
	}
	repo := NewBaseRepositoryImpl[entity.Example](relationField)
	return &ExampleSQLRepo{
		repo,
	}
}
