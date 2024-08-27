package main

import (
	"boiler-plate-clean/config"
	"boiler-plate-clean/internal/delivery/messaging"
	"boiler-plate-clean/internal/repository"
	services "boiler-plate-clean/internal/services"
	"context"
	kafkaserver "github.com/RumbiaID/pkg-library/app/pkg/broker/kafkaservice"
	"github.com/RumbiaID/pkg-library/app/pkg/database"
	"github.com/RumbiaID/pkg-library/app/pkg/logger"
	"github.com/RumbiaID/pkg-library/app/pkg/xvalidator"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func main() {
	validate, _ := xvalidator.NewValidator()
	conf := config.InitConsumerConfig(validate)
	logger.SetupLogger(&logger.Config{
		AppENV:  conf.AppEnvConfig.AppEnv,
		LogPath: conf.AppEnvConfig.LogFilePath,
		Debug:   conf.AppEnvConfig.AppDebug,
	})

	ctx, cancel := context.WithCancel(context.Background())
	// repository
	db := database.NewDatabase(conf.DatabaseConfig.Dbservice, &database.Config{
		DbHost:   conf.DatabaseConfig.Dbhost,
		DbUser:   conf.DatabaseConfig.Dbuser,
		DbPass:   conf.DatabaseConfig.Dbpassword,
		DbName:   conf.DatabaseConfig.Dbname,
		DbPort:   strconv.Itoa(conf.DatabaseConfig.Dbport),
		DbPrefix: conf.DatabaseConfig.DbPrefix,
	})
	userRepository := repository.NewUserRepository()
	walletRepository := repository.NewWalletRepository()
	categoryRepository := repository.NewCategoryTransactionRepository()
	transactionRepository := repository.NewTransactionRepository()
	fraudRepository := repository.NewFraudRepository()

	//service
	walletService := services.NewWalletService(db.GetDB(), walletRepository, userRepository, transactionRepository, categoryRepository, validate)
	transactionService := services.NewTransactionService(db.GetDB(), transactionRepository, categoryRepository, userRepository, walletRepository, nil, validate)
	fraudService := services.NewFraudService(db.GetDB(), fraudRepository, validate)

	//Handler
	transactionHandler := messaging.NewTransactionConsumer(transactionService, walletService, fraudService, conf.AppEnvConfig.LogFilePath)
	kafkaService := kafkaserver.New(&kafkaserver.Config{
		SecurityProtocol: conf.KafkaConfig.KafkaSecurityProtocol,
		Brokers:          conf.KafkaConfig.KafkaBroker,
		Username:         conf.KafkaConfig.KafkaUsername,
		Password:         conf.KafkaConfig.KafkaPassword,
	})
	go messaging.ConsumeKafkaTopic(ctx, kafkaService, conf.KafkaConfig.KafkaTopicTransaction, conf.KafkaConfig.KafkaGroupId, transactionHandler.ConsumeKafka)

	slog.Info("Worker is running")

	terminateSignals := make(chan os.Signal, 1)
	signal.Notify(terminateSignals, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM)

	stop := false
	for !stop {
		select {
		case s := <-terminateSignals:
			slog.Info("Got one of stop signals, shutting down worker gracefully, SIGNAL NAME :", s)
			cancel()
			stop = true
		}
	}

	time.Sleep(5 * time.Second) // wait for all consumers to finish processing
}
