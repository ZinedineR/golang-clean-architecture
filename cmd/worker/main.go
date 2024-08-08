package main

import (
	"boiler-plate-clean/config"
	"boiler-plate-clean/internal/delivery/messaging"
	kafkaMessaging "boiler-plate-clean/internal/delivery/messaging/kafka"
	redisMessaging "boiler-plate-clean/internal/delivery/messaging/redis"
	queueConsumer "boiler-plate-clean/internal/delivery/queue"
	"boiler-plate-clean/internal/gateway/queue"
	"context"
	"fmt"
	kafkaserver "github.com/RumbiaID/pkg-library/app/pkg/broker/kafkaservice"
	"github.com/RumbiaID/pkg-library/app/pkg/logger"
	"github.com/RumbiaID/pkg-library/app/pkg/redisser"
	"github.com/RumbiaID/pkg-library/app/pkg/twiliotext"
	"github.com/RumbiaID/pkg-library/app/pkg/xvalidator"
	"github.com/go-redis/redis/v8"
	amqp "github.com/rabbitmq/amqp091-go"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	rabbitClient *amqp.Connection
	kafkaService *kafkaserver.KafkaService
	redisClient  *redis.Client
	twilio       twiliotext.TwilioMethod
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
	for {
		conn, err := amqp.Dial(conf.RabbiterConfig.RabbitMQDial)

		if err != nil {
			slog.Error(fmt.Sprintf("failed to connect to rabbitmq database"), "error", err.Error())
			slog.Info(fmt.Sprintf("retrying to connect to rabbitmq in 5 seconds..."))
			time.Sleep(5 * time.Second)
			continue
		}
		slog.Info(fmt.Sprintf("successfully connected to rabbitmq"))
		rabbitClient = conn
		break
	}
	// repository

	// external api
	//httpClientFactory := httpclient.New()
	//httpClient := httpClientFactory.CreateClient()

	// queueproducer
	exampleProducer := queue.NewExampleProducerImpl(rabbitClient, "ExampleQueue")

	//Handler
	exampleHandler := messaging.NewExampleConsumer(exampleProducer)
	if conf.UsesKafka() {
		kafkaService = kafkaserver.New(&kafkaserver.Config{
			SecurityProtocol: conf.KafkaConfig.KafkaSecurityProtocol,
			Brokers:          conf.KafkaConfig.KafkaBroker,
			Username:         conf.KafkaConfig.KafkaUsername,
			Password:         conf.KafkaConfig.KafkaPassword,
		})
		go kafkaMessaging.ConsumeKafkaTopic(ctx, kafkaService, conf.KafkaConfig.KafkaTopicNotification, conf.KafkaConfig.KafkaGroupId, exampleHandler.ConsumeKafka)
	} else if conf.UsesRedis() {
		redisClient = redisser.NewRedis(&redisser.Config{
			Redisdatabase: conf.RedisConfig.Redisdatabase,
			Redishost:     conf.RedisConfig.Redishost,
			Redisport:     conf.RedisConfig.Redisport,
			Redispassword: conf.RedisConfig.Redispassword,
		})
		go redisMessaging.ConsumeRedisPublisher(ctx, redisClient, conf.AppName()+"-notification", exampleHandler.ConsumeRedis)
	}
	//queueConsumer
	exampleQueueHandler := queueConsumer.NewExampleConsumer()
	go queueConsumer.ConsumePublisher(ctx, rabbitClient, "ExampleQueue", exampleQueueHandler.Consume)

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
