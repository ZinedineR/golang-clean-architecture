package service

import (
	"boiler-plate-clean/pkg/exception"
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
)

type FileServiceImpl struct {
	filepath string
}

func NewFileService(filepath string) FileService {
	return &FileServiceImpl{filepath: filepath}
}

// CreateExample creates a new campaign
func (s *FileServiceImpl) Download(ctx context.Context) ([]byte, *exception.Exception) {
	file, err := os.OpenFile(s.filepath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return nil, exception.Internal("error opening file", err)
	}
	// Get the file size
	stat, err := file.Stat()
	if err != nil {
		fmt.Println(err)
		return nil, exception.Internal("error opening file", err)
	}

	// Read the file into a byte slice
	bs := make([]byte, stat.Size())
	_, err = bufio.NewReader(file).Read(bs)
	if err != nil && err != io.EOF {
		fmt.Println(err)
		return nil, exception.Internal("error reading file", err)
	}
	return bs, nil
}
