package setup

import (
	"context"
	"log"
	"log/slog"
	"strconv"
	"sync"
	"time"

	mq "github.com/dherik/ddd-golang-project/internal/infrastructure/messaging/rabbitmq" //FIXME
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	setupBrokerOnce    sync.Once
	RabbitMQDataSource mq.RabbitMQDataSource
	brokerContainer    *rabbitmq.RabbitMQContainer
	brokerCtx          context.Context
)

func SetupRabbitMQ() {
	setupBrokerOnce.Do(func() {
		slog.Info("Initializing the RabbitMQ...")

		ctx := context.Background()
		rabbitmqContainer, err := rabbitmq.RunContainer(ctx,
			testcontainers.WithImage("rabbitmq:3-management-alpine"),
			rabbitmq.WithAdminPassword("guest"),
			rabbitmq.WithAdminPassword("guest"),
			testcontainers.WithAfterReadyCommand(
				Exchange{
					Name:       "calendar",
					Type:       "fanout",
					AutoDelete: false,
					Durable:    true,
					Internal:   false,
				},
				Exchange{
					Name:       "billing",
					Type:       "direct",
					AutoDelete: false,
					Durable:    true,
					Internal:   false,
				},
			),
			testcontainers.WithWaitStrategy(
				wait.ForLog("Server startup complete").
					WithStartupTimeout(30*time.Second)),
		)
		if err != nil {
			log.Fatalf("failed to start RabbitMQ container: %s", err)
		}

		brokerContainer = rabbitmqContainer
		brokerCtx = ctx

		amqpUrl, err := rabbitmqContainer.AmqpURL(ctx)
		if err != nil {
			log.Fatalf("failed to get AMQP URL: %s", err)
		}
		slog.Debug("RabbitMQ connection URL", "url", amqpUrl)

		_, err = rabbitmqContainer.State(ctx)
		if err != nil {
			log.Fatalf("failed to get container state: %s", err)
		}

		portNat, _ := rabbitmqContainer.MappedPort(ctx, "5672") //FIXME
		port, _ := strconv.Atoi(portNat.Port())
		slog.Info("RabbitMQ for Integration Test is running", "port", port)

		host, _ := rabbitmqContainer.Host(ctx)

		rabbitmqDataSource := mq.RabbitMQDataSource{
			Host:     host,
			Port:     port,
			User:     rabbitmqContainer.AdminUsername,
			Password: rabbitmqContainer.AdminPassword,
			AmqpURL:  amqpUrl,
		}

		RabbitMQDataSource = rabbitmqDataSource
	})
}
