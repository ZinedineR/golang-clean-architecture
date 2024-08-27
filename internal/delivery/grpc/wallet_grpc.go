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

type WalletGRPCHandler struct {
	wallet_finance.UnimplementedWalletServiceServer
	WalletService services.WalletService
	GRPCParamHandler
}

func NewWalletGRPCHandler(service services.WalletService) *WalletGRPCHandler {
	return &WalletGRPCHandler{WalletService: service}
}

func (h *WalletGRPCHandler) CreateWallet(
	ctx context.Context, in *wallet_finance.CreateWalletRequest,
) (*wallet_finance.CreateWalletResponse, error) {
	wallet := &entity.Wallet{
		Name:   in.GetName(),
		UserId: int(in.GetUserId()),
	}
	if err := h.WalletService.Create(ctx, wallet); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.CreateWalletResponse{
		Wallet: &wallet_finance.Wallet{
			Id:              int32(wallet.ID),
			Name:            wallet.Name,
			UserId:          int32(wallet.UserId),
			Balance:         wallet.Balance,
			LastTransaction: timestamppb.New(*wallet.LastTransaction),
		},
		Response: &wallet_finance.MutationResponse{Message: "Create Wallet Success"},
	}, nil
}

func (h *WalletGRPCHandler) UpdateWallet(
	ctx context.Context, in *wallet_finance.UpdateWalletRequest,
) (*wallet_finance.UpdateWalletResponse, error) {
	wallet := &entity.Wallet{
		ID:     int(in.GetId()),
		Name:   in.GetName(),
		UserId: int(in.GetUserId()),
	}
	if err := h.WalletService.Update(ctx, wallet.ID, wallet); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.UpdateWalletResponse{
		Wallet: &wallet_finance.Wallet{
			Id:              int32(wallet.ID),
			Name:            wallet.Name,
			UserId:          int32(wallet.UserId),
			Balance:         wallet.Balance,
			LastTransaction: timestamppb.New(*wallet.LastTransaction),
		},
		Response: &wallet_finance.MutationResponse{Message: "Update Wallet Success"},
	}, nil
}

func (h *WalletGRPCHandler) GetWallet(
	ctx context.Context, in *wallet_finance.GetWalletRequest,
) (*wallet_finance.GetWalletResponse, error) {
	from, to, errParse := h.ParseDateParam(in.GetFromDate(), in.GetToDate())
	if errParse != nil {
		return nil, status.Error(codes.InvalidArgument, errParse.Error())
	}
	wallet, err := h.WalletService.DetailWalletTransaction(ctx, int(in.GetId()), from, to)
	if err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.GetWalletResponse{
		Wallet: &wallet_finance.Wallet{
			Id:              int32(wallet.ID),
			UserId:          int32(wallet.UserId),
			Balance:         wallet.Balance,
			LastTransaction: timestamppb.New(*wallet.LastTransaction),
		},
		Transactions: mapWalletTransactions(wallet.Transaction),
	}, nil
}

func (h *WalletGRPCHandler) GetLast10Transactions(
	ctx context.Context, in *wallet_finance.GetLast10TransactionsRequest,
) (*wallet_finance.GetLast10TransactionsResponse, error) {
	wallet, err := h.WalletService.Last10(ctx, int(in.GetId()))
	if err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.GetLast10TransactionsResponse{
		Wallet: &wallet_finance.Wallet{
			Id:              int32(wallet.ID),
			UserId:          int32(wallet.UserId),
			Balance:         wallet.Balance,
			LastTransaction: timestamppb.New(*wallet.LastTransaction),
		},
		Transactions: mapWalletTransactions(wallet.Transaction),
	}, nil
}

func (h *WalletGRPCHandler) RecapCategory(
	ctx context.Context, in *wallet_finance.RecapCategoryRequest,
) (*wallet_finance.RecapCategoryResponse, error) {
	from, to, errParse := h.ParseDateParam(in.GetFromDate(), in.GetToDate())
	if errParse != nil {
		return nil, status.Error(codes.InvalidArgument, errParse.Error())
	}
	walletRecap, err := h.WalletService.RecapCategory(ctx, int(in.GetId()), int(in.GetCategoryId()), from, to)
	if err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.RecapCategoryResponse{
		Id:         int32(walletRecap.ID),
		UserId:     int32(walletRecap.UserId),
		Total:      walletRecap.Total,
		CategoryId: int32(walletRecap.CategoryId),
		Category: &wallet_finance.CategoryTransaction{
			Id:        int32(walletRecap.Category.ID),
			Name:      walletRecap.Category.Name,
			CreatedAt: timestamppb.New(*walletRecap.Category.CreatedAt),
			UpdatedAt: timestamppb.New(*walletRecap.Category.UpdatedAt),
		},
		Transactions: mapWalletTransactions(walletRecap.Transaction),
	}, nil
}

func (h *WalletGRPCHandler) DeleteWallet(
	ctx context.Context, in *wallet_finance.DeleteWalletRequest,
) (*wallet_finance.DeleteWalletResponse, error) {
	if err := h.WalletService.Delete(ctx, int(in.GetId())); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.DeleteWalletResponse{
		Response: &wallet_finance.MutationResponse{Message: "Delete Wallet Success"},
	}, nil
}

// Helper functions to map entities to proto messages
func mapWalletTransactions(transactions []entity.Transaction) []*wallet_finance.Transaction {
	var protoTransactions []*wallet_finance.Transaction
	for _, transaction := range transactions {
		protoTransactions = append(protoTransactions, &wallet_finance.Transaction{
			Id:              int32(transaction.ID),
			Type:            transaction.Type,
			Amount:          transaction.Amount,
			Description:     transaction.Description,
			TransactionTime: timestamppb.New(*transaction.TransactionTime),
			WalletId:        int32(transaction.WalletId),
			CategoryId:      int32(transaction.CategoryId),
		})
	}
	return protoTransactions
}
