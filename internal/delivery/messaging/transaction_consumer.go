package messaging

import (
	"boiler-plate-clean/internal/entity"
	service "boiler-plate-clean/internal/services"
	"boiler-plate-clean/pkg/csvwriter"
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
	"log/slog"
	"os"
	"strings"
)

type TransactionConsumer struct {
	TransactionService service.TransactionService
	WalletService      service.WalletService
	FraudService       service.FraudService
	CSVWriter          csvwriter.CSVWriter[*entity.Transaction]
}

func NewTransactionConsumer(
	transactionService service.TransactionService,
	walletService service.WalletService,
	fraudService service.FraudService,
	filepath string,
) *TransactionConsumer {
	// Initialize the CSVWriter
	csvWriter, err := csvwriter.NewCSVWriter[*entity.Transaction](filepath, "/transaction.csv")
	if err != nil {
		slog.Error("error while creating csv writer", err)
		os.Exit(1)
	}
	if err := csvWriter.WriteToCSVHeader(&entity.Transaction{}); err != nil {
		slog.Error("error while creating csv writer", err)
		os.Exit(1)
	}
	return &TransactionConsumer{
		TransactionService: transactionService,
		WalletService:      walletService,
		FraudService:       fraudService,
		CSVWriter:          *csvWriter,
	}
}

func (c TransactionConsumer) ConsumeKafka(ctx context.Context, message *kafka.Message) error {
	transactionEvent := new(entity.Transaction)
	if err := json.Unmarshal(message.Value, transactionEvent); err != nil {
		slog.Error("error unmarshalling example event", slog.String("error", err.Error()))
		return err
	}
	//if err := c.CSVWriter.WriteToCSV(transactionEvent); err != nil {
	//	slog.Error("error writing example event", slog.String("error", err.Error()))
	//	return err
	//}

	if transactionEvent.Amount >= 1000000 && transactionEvent.Type == "transfer" {
		if err := c.FraudService.Create(ctx, &entity.Fraud{
			WalletId: transactionEvent.WalletId,
		}); err != nil {
			slog.Error("error creating fraud log", slog.String("error", err.Error.Error()))
			return err.Error
		}
		if err := c.TransactionService.Delete(ctx, transactionEvent.ID); err != nil {
			slog.Error("error deleting transaction", slog.String("error", err.Error.Error()))
			return err.Error
		}
		wallet, err := c.WalletService.Detail(ctx, transactionEvent.WalletId)
		if err != nil {
			slog.Error("error getting wallet")
			return err.Error
		}
		walletUpdate := entity.Wallet{
			ID:      transactionEvent.WalletId,
			Name:    wallet.Name,
			UserId:  wallet.UserId,
			Balance: wallet.Balance,
		}
		if strings.Contains(transactionEvent.Description, "from") {
			walletUpdate.Decrease(transactionEvent.Amount)
			fmt.Println(transactionEvent.Description)
			fmt.Println("previous balance:", wallet.Balance)
			fmt.Println("current balance:", walletUpdate.Balance)
		} else if strings.Contains(transactionEvent.Description, "to") {
			walletUpdate.Increase(transactionEvent.Amount)
			fmt.Println(transactionEvent.Description)
			fmt.Println("previous balance:", wallet.Balance)
			fmt.Println("current balance:", walletUpdate.Balance)
		}
		if err := c.WalletService.Update(ctx, transactionEvent.WalletId, &walletUpdate); err != nil {
			slog.Error("error updating wallet")
			return err.Error
		}
	}

	slog.Info("Received topic example with event", slog.Any("example", transactionEvent))
	return nil
}
