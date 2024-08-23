package main

import (
	"boiler-plate-clean/config"
	"boiler-plate-clean/internal/delivery/messaging"
	"boiler-plate-clean/internal/repository"
	service "boiler-plate-clean/internal/services"
	"boiler-plate-clean/migration"
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

var (
	kafkaService *kafkaserver.KafkaService
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
	//ctx, cancel := context.WithCancel(context.Background())
	// db
	userRepository := repository.NewUserRepository()
	userWriteDB := initSQLWrite(conf)
	userReadDB := initSQLRead(conf)

	// service
	userWriteService := service.NewUserService(userWriteDB.GetDB(), userRepository, validate)
	userReadService := service.NewUserService(userReadDB.GetDB(), userRepository, validate)
	//httpClient := httpClientFactory.CreateClient()

	//Handler
	userWriteHandler := messaging.NewUserWriteConsumer(userWriteService, userReadService)

	kafkaService = kafkaserver.New(&kafkaserver.Config{
		SecurityProtocol: conf.KafkaConfig.KafkaSecurityProtocol,
		Brokers:          conf.KafkaConfig.KafkaBroker,
		Username:         conf.KafkaConfig.KafkaUsername,
		Password:         conf.KafkaConfig.KafkaPassword,
	})
	go messaging.ConsumeKafkaTopic(ctx, kafkaService, conf.KafkaConfig.KafkaTopic, conf.KafkaConfig.KafkaGroupId, userWriteHandler.ConsumeKafka)
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

func initSQLWrite(conf *config.Config) *database.Database {
	db := database.NewDatabase(conf.DatabaseConfig.Dbservice, &database.Config{
		DbHost:   conf.DatabaseConfig.Dbhost,
		DbUser:   conf.DatabaseConfig.Dbuser,
		DbPass:   conf.DatabaseConfig.Dbpassword,
		DbName:   conf.DatabaseConfig.Dbname,
		DbPort:   strconv.Itoa(conf.DatabaseConfig.Dbport),
		DbPrefix: conf.DatabaseConfig.DbPrefix,
	})
	if conf.IsStaging() {
		migration.AutoMigration(db)
	}
	return db
}
func initSQLRead(conf *config.Config) *database.Database {
	db := database.NewDatabase(conf.DatabaseConfig.Dbservice, &database.Config{
		DbHost:   conf.DatabaseReplicaConfig.Dbreplicahost,
		DbUser:   conf.DatabaseReplicaConfig.Dbreplicauser,
		DbPass:   conf.DatabaseReplicaConfig.Dbreplicapassword,
		DbName:   conf.DatabaseReplicaConfig.Dbreplicaname,
		DbPort:   strconv.Itoa(conf.DatabaseReplicaConfig.Dbreplicaport),
		DbPrefix: conf.DatabaseConfig.DbPrefix,
	})
	if conf.IsStaging() {
		migration.AutoMigration(db)
	}
	return db
}
