package grpc

import (
	"boiler-plate-clean/internal/model"
	"fmt"
	"google.golang.org/protobuf/types/known/timestamppb"
	"regexp"
	"strings"
	"time"
)

const (
	invalidParameter = "invalid %s parameter"
)

type GRPCParamHandler struct {
}

func (h *GRPCParamHandler) ParseDateParam(from, to *timestamppb.Timestamp) (time.Time, time.Time, error) {
	toDate := time.Now().AddDate(0, 0, 1)
	var fromDate time.Time
	var err error
	if !from.AsTime().IsZero() {
		fromDate, err = time.Parse("2006-01-02", from.AsTime().Format("2006-01-02"))
		if err != nil {
			return fromDate, toDate, err
		}
	}
	if !to.AsTime().IsZero() {
		toDate, err = time.Parse("2006-01-02", to.AsTime().Format("2006-01-02"))
		if err != nil {
			return fromDate, toDate, err
		}
		toDate = toDate.AddDate(0, 0, 1)
	}
	return fromDate, toDate, nil
}

var orderRegex = regexp.MustCompile("(\\w+):(\\w+)")

var OrderOperators = map[string]string{
	"desc": "desc",
	"asc":  "asc",
}

func GetOrderValue(value string) (string, error) {
	if op, ok := OrderOperators[value]; ok {
		return op, nil
	}
	return "", fmt.Errorf(invalidParameter, value)
}

var filterRegex = regexp.MustCompile(`(\w+):([^|]+):(\w+)`)

var FilterOperator = map[string]string{
	"eq":   "=",
	"lt":   "<",
	"gt":   ">",
	"lte":  "<=",
	"gte":  ">=",
	"in":   "in",
	"like": "like",
	"is":   "is",
	"not":  "not in",
}

func GetFilterOperator(operator string) (string, error) {
	if op, ok := FilterOperator[operator]; ok {
		return op, nil
	}
	return "", fmt.Errorf(invalidParameter, operator)
}

func (h *GRPCParamHandler) ParseOrderParam(order string) (model.OrderParam, error) {
	var p model.OrderParam
	if order != "" {
		listOrder := strings.Split(order, ",")
		for _, o := range listOrder {
			if !orderRegex.MatchString(o) {
				continue
			}
			condition := strings.Split(o, ":")
			if len(condition) != 2 {
				return model.OrderParam{}, fmt.Errorf(invalidParameter, "order")
			}
			value, err := GetOrderValue(condition[1])
			if err != nil {
				return model.OrderParam{}, err
			}
			p.OrderBy = condition[0]
			p.Order = value
		}
	}
	return p, nil
}

func (h *GRPCParamHandler) ParseFilterParams(f string) (model.FilterParams, error) {
	var p model.FilterParams

	if f != "" {
		listFilter := strings.Split(f, "|")
		for _, v := range listFilter {
			if !filterRegex.MatchString(v) {
				continue
			}
			filter := strings.Split(v, ":")
			if len(filter) != 3 {
				return model.FilterParams{}, fmt.Errorf(invalidParameter, filter)
			}
			operator, err := GetFilterOperator(filter[2])
			if err != nil {
				return model.FilterParams{}, err
			}
			p = append(p, &model.FilterParam{
				Field:    filter[0],
				Value:    filter[1],
				Operator: operator,
			})
		}
	}

	return p, nil
}

func (h *GRPCParamHandler) ParseFindParams(order, filter string) (
	model.OrderParam, model.FilterParams, error,
) {
	orders, err := h.ParseOrderParam(order)
	if err != nil {
		return model.OrderParam{}, model.FilterParams{}, err
	}
	filters, err := h.ParseFilterParams(filter)
	if err != nil {
		return model.OrderParam{}, model.FilterParams{}, err
	}
	return orders, filters, nil
}
