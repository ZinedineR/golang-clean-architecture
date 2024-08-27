package csvwriter

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
)

// CSVWriter is a generic struct that can write any type of data to a CSV file.
type CSVWriter[T any] struct {
	FilePath string
	Writer   *csv.Writer
}

// NewCSVWriter creates a new CSVWriter instance, ensuring the directory exists, and initializes the CSV writer.
func NewCSVWriter[T any](dir string, filename string) (*CSVWriter[T], error) {
	// Ensure the directory exists
	err := os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return nil, fmt.Errorf("failed to create directory: %w", err)
	}

	// Create or open the CSV file
	filePath := filepath.Join(dir, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create file: %w", err)
	}

	// Initialize CSV writer
	writer := csv.NewWriter(file)

	return &CSVWriter[T]{
		FilePath: filePath,
		Writer:   writer,
	}, nil
}

// WriteToCSV writes a slice of generic type T to the CSV file.
func (cw *CSVWriter[T]) WriteToCSV(data ...T) error {
	for _, record := range data {
		// Convert the struct to a slice of strings
		recordSlice := structToSlice(record)
		err := cw.Writer.Write(recordSlice)
		if err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	// Flush to ensure data is written to file
	cw.Writer.Flush()
	if err := cw.Writer.Error(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

func (cw *CSVWriter[T]) WriteToCSVHeader(data ...T) error {
	file, err := os.Open(cw.FilePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	_, err = reader.Read()
	if err == nil {
		// Data already exists in the file, return nil
		return nil
	}
	for _, record := range data {
		// Convert the struct to a slice of strings
		recordSlice := structToSliceHeader(record)
		err := cw.Writer.Write(recordSlice)
		if err != nil {
			return fmt.Errorf("failed to write record: %w", err)
		}
	}

	// Flush to ensure data is written to file
	cw.Writer.Flush()
	if err := cw.Writer.Error(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}

	return nil
}

// Close closes the underlying CSV writer and file.
func (cw *CSVWriter[T]) Close() error {
	cw.Writer.Flush()
	if err := cw.Writer.Error(); err != nil {
		return fmt.Errorf("failed to flush writer: %w", err)
	}
	return nil
}

// structToSlice converts a struct to a slice of strings for writing to CSV.
func structToSlice[T any](data T) []string {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	var result []string
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		var value string

		switch field.Kind() {
		case reflect.String:
			value = field.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = strconv.FormatInt(field.Int(), 10)
		case reflect.Float32, reflect.Float64:
			value = strconv.FormatFloat(field.Float(), 'f', -1, 64)
		case reflect.Bool:
			value = strconv.FormatBool(field.Bool())
		default:
			continue
		}

		result = append(result, value)
	}

	return result
}

func structToSliceHeader[T any](data T) []string {
	body := reflect.ValueOf(data)
	if body.Kind() == reflect.Ptr {
		body = body.Elem()
	}
	bodyType := reflect.TypeOf(data)
	if bodyType.Kind() == reflect.Ptr {
		bodyType = bodyType.Elem()
	}

	var result []string
	for i := 0; i < body.NumField(); i++ {
		field := body.Field(i)
		fieldType := bodyType.Field(i)
		jsonTag := fieldType.Tag.Get("json")
		var value string

		switch field.Kind() {
		case reflect.String:
			value = jsonTag
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			value = jsonTag
		case reflect.Float32, reflect.Float64:
			value = jsonTag
		case reflect.Bool:
			value = jsonTag
		default:
			continue
		}

		result = append(result, value)
	}

	return result
}
