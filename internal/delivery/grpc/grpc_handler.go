package grpc

type BaseGRPCHandler struct {
	Category    *CategoryTransactionGRPCHandler
	Transaction *TransactionGRPCHandler
	User        *UserGRPCHandler
	Wallet      *WalletGRPCHandler
}

func NewBaseGRPCHandler(
	category *CategoryTransactionGRPCHandler,
	transaction *TransactionGRPCHandler,
	user *UserGRPCHandler,
	wallet *WalletGRPCHandler,
) *BaseGRPCHandler {
	return &BaseGRPCHandler{
		Category:    category,
		Transaction: transaction,
		User:        user,
		Wallet:      wallet,
	}
}
