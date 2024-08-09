package grpc

import (
	"boiler-plate-clean/internal/entity"
	"boiler-plate-clean/internal/model"
	service "boiler-plate-clean/internal/services"
	wallet_finance "boiler-plate-clean/proto/wallet-finance/v1"
	"context"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type UserGRPCHandler struct {
	wallet_finance.UnimplementedUserServiceServer
	UserService service.UserService
}

func NewUserGRPCHandler(service service.UserService) *UserGRPCHandler {
	return &UserGRPCHandler{UserService: service}
}

func (h *UserGRPCHandler) CreateUser(
	ctx context.Context, in *wallet_finance.CreateUserRequest,
) (*wallet_finance.CreateUserResponse, error) {
	user := &entity.Users{
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}
	if err := h.UserService.Create(ctx, user); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.CreateUserResponse{
		User: &wallet_finance.Users{
			Id:        int32(user.ID),
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: timestamppb.New(*user.CreatedAt),
			UpdatedAt: timestamppb.New(*user.UpdatedAt),
		},
		Response: &wallet_finance.MutationResponse{Message: "Create User Success"},
	}, nil
}

func (h *UserGRPCHandler) UpdateUser(
	ctx context.Context, in *wallet_finance.UpdateUserRequest,
) (*wallet_finance.UpdateUserResponse, error) {
	user := &entity.Users{
		ID:       int(in.GetId()),
		Email:    in.GetEmail(),
		Password: in.GetPassword(),
	}
	if err := h.UserService.Update(ctx, user.ID, user); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.UpdateUserResponse{
		User: &wallet_finance.Users{
			Id:        int32(user.ID),
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: timestamppb.New(*user.CreatedAt),
			UpdatedAt: timestamppb.New(*user.UpdatedAt),
		},
		Response: &wallet_finance.MutationResponse{Message: "Update User Success"},
	}, nil
}

func (h *UserGRPCHandler) DetailUser(
	ctx context.Context, in *wallet_finance.DetailUserRequest,
) (*wallet_finance.DetailUserResponse, error) {
	user, err := h.UserService.Detail(ctx, int(in.GetId()))
	if err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.DetailUserResponse{
		User: &wallet_finance.Users{
			Id:        int32(user.ID),
			Email:     user.Email,
			Password:  user.Password,
			CreatedAt: timestamppb.New(*user.CreatedAt),
			UpdatedAt: timestamppb.New(*user.UpdatedAt),
			Wallets:   mapWallets(user.Wallet),
		},
	}, nil
}

func (h *UserGRPCHandler) DeleteUser(
	ctx context.Context, in *wallet_finance.DeleteUserRequest,
) (*wallet_finance.DeleteUserResponse, error) {
	if err := h.UserService.Delete(ctx, int(in.GetId())); err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.DeleteUserResponse{
		Response: &wallet_finance.MutationResponse{Message: "Delete User Success"},
	}, nil
}

func (h *UserGRPCHandler) Cashflow(
	ctx context.Context, in *wallet_finance.CashflowRequest,
) (*wallet_finance.CashflowResponse, error) {
	userWalletRecap, err := h.UserService.Cashflow(ctx, int(in.GetId()), in.GetFromDate().AsTime(), in.GetToDate().AsTime())
	if err != nil {
		return nil, status.Error(codes.Code(err.GetGrpcCode()), fmt.Sprint(err.Message))
	}
	return &wallet_finance.CashflowResponse{
		Id:        int32(userWalletRecap.ID),
		Email:     userWalletRecap.Email,
		CreatedAt: timestamppb.New(*userWalletRecap.CreatedAt),
		UpdatedAt: timestamppb.New(*userWalletRecap.UpdatedAt),
		Wallets:   mapWalletDetails(userWalletRecap.Wallet),
	}, nil
}

// Helper functions to map entities to proto messages
func mapWallets(wallets []entity.Wallet) []*wallet_finance.Wallet {
	var protoWallets []*wallet_finance.Wallet
	for _, wallet := range wallets {
		protoWallets = append(protoWallets, &wallet_finance.Wallet{
			Id:              int32(wallet.ID),
			Name:            wallet.Name,
			Balance:         wallet.Balance,
			LastTransaction: timestamppb.New(*wallet.LastTransaction),
		})
	}
	return protoWallets
}

func mapWalletDetails(walletDetails []model.UserWalletDetail) []*wallet_finance.WalletDetail {
	var protoWalletDetails []*wallet_finance.WalletDetail
	for _, walletDetail := range walletDetails {
		protoWalletDetails = append(protoWalletDetails, &wallet_finance.WalletDetail{
			Id:               int32(walletDetail.ID),
			Name:             walletDetail.Name,
			TotalIncome:      walletDetail.TotalIncome,
			TotalExpense:     walletDetail.TotalExpense,
			WalletTypeDetail: mapWalletTypeDetails(walletDetail.UserWalletTypeDetail),
		})
	}
	return protoWalletDetails
}

func mapWalletTypeDetails(walletTypeDetails []model.UserWalletTypeDetail) []*wallet_finance.WalletTypeDetail {
	var protoWalletTypeDetails []*wallet_finance.WalletTypeDetail
	for _, walletTypeDetail := range walletTypeDetails {
		protoWalletTypeDetails = append(protoWalletTypeDetails, &wallet_finance.WalletTypeDetail{
			Type:         walletTypeDetail.Type,
			Total:        walletTypeDetail.Total,
			Transactions: mapTransactions(walletTypeDetail.Transaction),
		})
	}
	return protoWalletTypeDetails
}

func mapTransactions(transactions []entity.Transaction) []*wallet_finance.Transaction {
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
