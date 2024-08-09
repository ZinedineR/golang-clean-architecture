package route

import (
	"boiler-plate-clean/pkg/server"
	wallet_finance "boiler-plate-clean/proto/wallet-finance/v1"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"log"
)

func (h *Router) GRPCSetup() {
	grpcServer, netListener := server.NewGRPCServer("50051")
	grpcGatewayMux := runtime.NewServeMux()

	wallet_finance.RegisterUserServiceServer(grpcServer, h.GRPCHandler.User)
	wallet_finance.RegisterWalletServiceServer(grpcServer, h.GRPCHandler.Wallet)
	wallet_finance.RegisterCategoryTransactionServiceServer(grpcServer, h.GRPCHandler.Category)
	wallet_finance.RegisterTransactionServiceServer(grpcServer, h.GRPCHandler.Transaction)

	grpcClient := server.NewGRPCClient(grpcServer, netListener, "50051")

	if err := wallet_finance.RegisterUserServiceHandler(context.Background(), grpcGatewayMux, grpcClient); err != nil {
		log.Fatal("failed to register users gateway")
	}
	if err := wallet_finance.RegisterWalletServiceHandler(context.Background(), grpcGatewayMux, grpcClient); err != nil {
		log.Fatal("failed to register wallets gateway")
	}
	if err := wallet_finance.RegisterCategoryTransactionServiceHandler(context.Background(), grpcGatewayMux, grpcClient); err != nil {
		log.Fatal("failed to register category transactions gateway")
	}
	if err := wallet_finance.RegisterTransactionServiceHandler(context.Background(), grpcGatewayMux, grpcClient); err != nil {
		log.Fatal("failed to register transactions gateway")
	}
	h.App.Group("/api/v2/*{grpc_gateway}").Any("", gin.WrapH(grpcGatewayMux))
}
