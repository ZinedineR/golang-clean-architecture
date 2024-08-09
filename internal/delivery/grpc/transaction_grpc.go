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

type TransactionGRPCHandler struct {
	wallet_finance.UnimplementedTransactionServiceServer
	TransactionService services.TransactionService
}

func NewTransactionGRPCHandler(service services.TransactionService) *TransactionGRPCHandler {
	return &TransactionGRPCHandler{TransactionService: service}
}

func (h *TransactionGRPCHandler) CreateTransaction(
	ctx context.Context, in *wallet_finance.CreateTransactionRequest,
) (*wallet_finance.CreateTransactionResponse, error) {
	transaction := &entity.Transaction{
		Type:        in.GetType(),
		Amount:      in.GetAmount(),
		Description: in.GetDescription(),
		WalletId:    int(in.GetWalletId()),
		CategoryId:  int(in.GetCategoryId()),
	}
	if err := h.TransactionService.Create(ctx, transaction); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.CreateTransactionResponse{
		Transaction: &wallet_finance.Transaction{
			Id:              int32(transaction.ID),
			Type:            transaction.Type,
			Amount:          transaction.Amount,
			Description:     transaction.Description,
			TransactionTime: timestamppb.New(*transaction.TransactionTime),
			WalletId:        int32(transaction.WalletId),
			CategoryId:      int32(transaction.CategoryId),
		},
		Response: &wallet_finance.MutationResponse{Message: "Create Transaction Success"},
	}, nil
}

func (h *TransactionGRPCHandler) UpdateTransaction(
	ctx context.Context, in *wallet_finance.UpdateTransactionRequest,
) (*wallet_finance.UpdateTransactionResponse, error) {
	transaction := &entity.Transaction{
		ID:          int(in.GetId()),
		Type:        in.GetType(),
		Amount:      in.GetAmount(),
		Description: in.GetDescription(),
		WalletId:    int(in.GetWalletId()),
		CategoryId:  int(in.GetCategoryId()),
	}
	if err := h.TransactionService.Update(ctx, transaction.ID, transaction); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.UpdateTransactionResponse{
		Transaction: &wallet_finance.Transaction{
			Id:              int32(transaction.ID),
			Type:            transaction.Type,
			Amount:          transaction.Amount,
			Description:     transaction.Description,
			TransactionTime: timestamppb.New(*transaction.TransactionTime),
			WalletId:        int32(transaction.WalletId),
			CategoryId:      int32(transaction.CategoryId),
		},
		Response: &wallet_finance.MutationResponse{Message: "Update Transaction Success"},
	}, nil
}

func (h *TransactionGRPCHandler) GetTransaction(
	ctx context.Context, in *wallet_finance.GetTransactionRequest,
) (*wallet_finance.GetTransactionResponse, error) {
	transaction, err := h.TransactionService.Detail(ctx, int(in.GetId()))
	if err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.GetTransactionResponse{
		Transaction: &wallet_finance.Transaction{
			Id:              int32(transaction.ID),
			Type:            transaction.Type,
			Amount:          transaction.Amount,
			Description:     transaction.Description,
			TransactionTime: timestamppb.New(*transaction.TransactionTime),
			WalletId:        int32(transaction.WalletId),
			CategoryId:      int32(transaction.CategoryId),
		},
	}, nil
}

func (h *TransactionGRPCHandler) CreditTransaction(
	ctx context.Context, in *wallet_finance.CreditTransactionRequest,
) (*wallet_finance.CreditTransactionResponse, error) {
	if err := h.TransactionService.Credit(ctx, int(in.GetWalletId()), int(in.GetCategoryId()), in.GetAmount()); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.CreditTransactionResponse{
		Response: &wallet_finance.MutationResponse{Message: "Credit Transaction Success"},
	}, nil
}

func (h *TransactionGRPCHandler) TransferTransaction(
	ctx context.Context, in *wallet_finance.TransferTransactionRequest,
) (*wallet_finance.TransferTransactionResponse, error) {
	if err := h.TransactionService.Transfer(ctx, int(in.GetSenderId()), int(in.GetReceiverId()), in.GetAmount()); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.TransferTransactionResponse{
		Response: &wallet_finance.MutationResponse{Message: "Transfer Transaction Success"},
	}, nil
}

func (h *TransactionGRPCHandler) DeleteTransaction(
	ctx context.Context, in *wallet_finance.DeleteTransactionRequest,
) (*wallet_finance.DeleteTransactionResponse, error) {
	if err := h.TransactionService.Delete(ctx, int(in.GetId())); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.DeleteTransactionResponse{
		Response: &wallet_finance.MutationResponse{Message: "Delete Transaction Success"},
	}, nil
}
