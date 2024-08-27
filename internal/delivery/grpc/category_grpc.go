package grpc

import (
	"boiler-plate-clean/internal/entity"
	services "boiler-plate-clean/internal/services"
	wallet_finance "boiler-plate-clean/proto/wallet-finance/v1"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type CategoryTransactionGRPCHandler struct {
	wallet_finance.UnimplementedCategoryTransactionServiceServer
	CategoryTransactionService services.CategoryTransactionService
	GRPCParamHandler
}

func NewCategoryTransactionGRPCHandler(service services.CategoryTransactionService) *CategoryTransactionGRPCHandler {
	return &CategoryTransactionGRPCHandler{CategoryTransactionService: service}
}

func (h *CategoryTransactionGRPCHandler) CreateCategoryTransaction(
	ctx context.Context, in *wallet_finance.CreateCategoryTransactionRequest,
) (*wallet_finance.CreateCategoryTransactionResponse, error) {
	categoryTransaction := &entity.CategoryTransaction{
		Name: in.GetName(),
	}
	if err := h.CategoryTransactionService.Create(ctx, categoryTransaction); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.CreateCategoryTransactionResponse{
		CategoryTransaction: &wallet_finance.CategoryTransaction{
			Id:        int32(categoryTransaction.ID),
			Name:      categoryTransaction.Name,
			CreatedAt: timestamppb.New(*categoryTransaction.CreatedAt),
			UpdatedAt: timestamppb.New(*categoryTransaction.UpdatedAt),
		},
		Response: &wallet_finance.MutationResponse{Message: "Create Category Transaction Success"},
	}, nil
}

func (h *CategoryTransactionGRPCHandler) UpdateCategoryTransaction(
	ctx context.Context, in *wallet_finance.UpdateCategoryTransactionRequest,
) (*wallet_finance.UpdateCategoryTransactionResponse, error) {
	categoryTransaction := &entity.CategoryTransaction{
		ID:   int(in.GetId()),
		Name: in.GetName(),
	}
	if err := h.CategoryTransactionService.Update(ctx, categoryTransaction.ID, categoryTransaction); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.UpdateCategoryTransactionResponse{
		CategoryTransaction: &wallet_finance.CategoryTransaction{
			Id:        int32(categoryTransaction.ID),
			Name:      categoryTransaction.Name,
			CreatedAt: timestamppb.New(*categoryTransaction.CreatedAt),
			UpdatedAt: timestamppb.New(*categoryTransaction.UpdatedAt),
		},
		Response: &wallet_finance.MutationResponse{Message: "Update Category Transaction Success"},
	}, nil
}

func (h *CategoryTransactionGRPCHandler) GetCategoryTransaction(
	ctx context.Context, in *wallet_finance.GetCategoryTransactionRequest,
) (*wallet_finance.GetCategoryTransactionResponse, error) {
	categoryTransaction, err := h.CategoryTransactionService.Detail(ctx, int(in.GetId()))
	if err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.GetCategoryTransactionResponse{
		CategoryTransaction: &wallet_finance.CategoryTransaction{
			Id:        int32(categoryTransaction.ID),
			Name:      categoryTransaction.Name,
			CreatedAt: timestamppb.New(*categoryTransaction.CreatedAt),
			UpdatedAt: timestamppb.New(*categoryTransaction.UpdatedAt),
		},
	}, nil
}

func (h *CategoryTransactionGRPCHandler) FindCategoryTransactions(
	ctx context.Context, in *wallet_finance.FindCategoryTransactionsRequest,
) (*wallet_finance.FindCategoryTransactionsResponse, error) {
	order, filter, errParam := h.ParseFindParams(in.GetOrder(), in.GetFilter())
	if errParam != nil {
		return nil, status.Error(codes.InvalidArgument, errParam.Error())
	}
	categoryTransactions, err := h.CategoryTransactionService.Find(ctx, order, filter)
	if err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	protoCategoryTransactions := []*wallet_finance.CategoryTransaction{}
	for _, ct := range *categoryTransactions {
		protoCategoryTransactions = append(protoCategoryTransactions, &wallet_finance.CategoryTransaction{
			Id:        int32(ct.ID),
			Name:      ct.Name,
			CreatedAt: timestamppb.New(*ct.CreatedAt),
			UpdatedAt: timestamppb.New(*ct.UpdatedAt),
		})
	}
	return &wallet_finance.FindCategoryTransactionsResponse{
		CategoryTransactions: protoCategoryTransactions,
	}, nil
}

func (h *CategoryTransactionGRPCHandler) DeleteCategoryTransaction(
	ctx context.Context, in *wallet_finance.DeleteCategoryTransactionRequest,
) (*wallet_finance.DeleteCategoryTransactionResponse, error) {
	if err := h.CategoryTransactionService.Delete(ctx, int(in.GetId())); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.DeleteCategoryTransactionResponse{
		Response: &wallet_finance.MutationResponse{Message: "Delete Category Transaction Success"},
	}, nil
}
