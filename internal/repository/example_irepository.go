package repository

import (
	"boiler-plate-clean/internal/entity"
)

type ExampleRepository interface {
	// Example operations
	BaseRepository[entity.Example]
}
