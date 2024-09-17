package service

import (
	"boiler-plate-clean/pkg/exception"
	"context"
)

type FileService interface {
	// CRUD operations for Example
	Download(ctx context.Context) ([]byte, *exception.Exception)
}
