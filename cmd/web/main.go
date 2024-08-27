package main

import (
	"boiler-plate-clean/config"
	"boiler-plate-clean/internal/delivery/grpc"
	"boiler-plate-clean/internal/delivery/http"
	"boiler-plate-clean/internal/delivery/http/route"
	"boiler-plate-clean/internal/gateway/messaging"
	"boiler-plate-clean/internal/repository"
	services "boiler-plate-clean/internal/services"
	"boiler-plate-clean/migration"
	"boiler-plate-clean/pkg/server"
	kafkaserver "github.com/RumbiaID/pkg-library/app/pkg/broker/kafkaservice"
	"github.com/RumbiaID/pkg-library/app/pkg/database"
	"github.com/RumbiaID/pkg-library/app/pkg/httpclient"
	"github.com/RumbiaID/pkg-library/app/pkg/logger"
	"github.com/RumbiaID/pkg-library/app/pkg/xvalidator"
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"
)

var (
	sqlClientRepo *database.Database
	kafkaDialer   *kafkaserver.KafkaService
)

// @title           Pigeon
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8080
// @BasePath  /api/notificationsvc/api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	validate, _ := xvalidator.NewValidator()
	conf := config.InitAppConfig(validate)
	logger.SetupLogger(&logger.Config{
		AppENV:  conf.AppEnvConfig.AppEnv,
		LogPath: conf.AppEnvConfig.LogFilePath,
		Debug:   conf.AppEnvConfig.AppDebug,
	})
	initInfrastructure(conf)
	ginServer := server.NewGinServer(&server.GinConfig{
		HttpPort:     conf.AppEnvConfig.HttpPort,
		AllowOrigins: conf.AppEnvConfig.AllowOrigins,
		AllowMethods: conf.AppEnvConfig.AllowMethods,
		AllowHeaders: conf.AppEnvConfig.AllowHeaders,
	})

	// repository
	userRepository := repository.NewUserRepository()
	walletRepository := repository.NewWalletRepository()
	categoryRepository := repository.NewCategoryTransactionRepository()
	transactionRepository := repository.NewTransactionRepository()

	// external api
	//gotifySvcExternalAPI := externalapi.NewExampleExternalImpl(conf, httpClient)

	// producer
	//if conf.UsesRedis() {
	//	exampleProducer = messaging.NewExampleRedisProducerImpl(redisClient, conf.AppName()+"-email")
	//} else if conf.UsesKafka() {
	//	exampleProducer = messaging.NewWalletProducerImpl(kafkaDialer, conf.KafkaConfig.KafkaTopicEmail)
	//}
	transactionProducer := messaging.NewTransactionProducerImpl(kafkaDialer, conf.KafkaConfig.KafkaTopicTransaction)

	// service
	userService := services.NewUserService(sqlClientRepo.GetDB(), userRepository, walletRepository, transactionRepository, validate)
	walletService := services.NewWalletService(sqlClientRepo.GetDB(), walletRepository, userRepository, transactionRepository, categoryRepository, validate)
	categoryService := services.NewCategoryTransactionService(sqlClientRepo.GetDB(), categoryRepository, validate)
	transactionService := services.NewTransactionService(sqlClientRepo.GetDB(), transactionRepository, categoryRepository, userRepository, walletRepository, transactionProducer, validate)
	// Handler
	userHandler := http.NewUserHTTPHandler(userService)
	walletHandler := http.NewWalletHTTPHandler(walletService)
	categoryHandler := http.NewCategoryTransactionHTTPHandler(categoryService)
	transactionHandler := http.NewTransactionHTTPHandler(transactionService)

	//GRPC
	userGRPC := grpc.NewUserGRPCHandler(userService)
	walletGRPC := grpc.NewWalletGRPCHandler(walletService)
	categoryGRPC := grpc.NewCategoryTransactionGRPCHandler(categoryService)
	transactionGRPC := grpc.NewTransactionGRPCHandler(transactionService)
	serverGRPC := grpc.NewBaseGRPCHandler(categoryGRPC, transactionGRPC, userGRPC, walletGRPC)

	router := route.Router{
		App:                ginServer.App,
		UserHandler:        userHandler,
		WalletHandler:      walletHandler,
		CategoryHandler:    categoryHandler,
		TransactionHandler: transactionHandler,
		Config:             conf,
		GRPCHandler:        serverGRPC,
	}
	//router.Setup()
	router.GRPCSetup()
	//router.SwaggerRouter()

	echan := make(chan error)
	go func() {
		echan <- ginServer.Start()
	}()

	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM)

	select {
	case <-term:
		slog.Info("signal terminated detected")
	case err := <-echan:
		slog.Error("Failed to start http server", err)
	}
}

func initInfrastructure(config *config.Config) {
	kafkaDialer = initKafka(config)
	sqlClientRepo = initSQL(config)
}

func initSQL(conf *config.Config) *database.Database {
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

func initHttpclient() httpclient.Client {
	httpClientFactory := httpclient.New()
	httpClient := httpClientFactory.CreateClient()
	return httpClient
}

func initKafka(config *config.Config) *kafkaserver.KafkaService {
	kafkaDialer := kafkaserver.New(&kafkaserver.Config{
		SecurityProtocol: config.KafkaConfig.KafkaSecurityProtocol,
		Brokers:          config.KafkaConfig.KafkaBroker,
		Username:         config.KafkaConfig.KafkaUsername,
		Password:         config.KafkaConfig.KafkaPassword,
	})
	return kafkaDialer
}
