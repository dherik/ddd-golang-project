package setup

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/dherik/ddd-golang-project/internal/infrastructure/messaging"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	initializedBroker     bool
	initializedBrokerLock sync.Mutex
	// PostgresContainer *postgres.PostgresContainer
	RabbitMQDataSource messaging.RabbitMQDataSource
	// Ctx               context.Context
)

func SetupRabbitMQ() {

	initializedBrokerLock.Lock()
	defer initializedBrokerLock.Unlock()

	if initializedBroker {
		slog.Debug("RabbitMQ it's already initialized")
		return
	}

	slog.Info("Initializing the RabbitMQ...")

	ctx := context.Background()
	rabbitmqContainer, err := rabbitmq.RunContainer(ctx,
		testcontainers.WithImage("rabbitmq:3-management-alpine"),
		rabbitmq.WithAdminPassword("guest"),
		rabbitmq.WithAdminPassword("guest"),
		testcontainers.WithAfterReadyCommand(Exchange{
			Name:       "calendar",
			Type:       "fanout",
			AutoDelete: false,
			Durable:    true,
			Internal:   false,
		}),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Server startup complete").
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start RabbitMQ container: %s", err)
	}

	amqpUrl, err := rabbitmqContainer.AmqpURL(ctx)
	if err != nil {
		log.Fatalf("failed to get AMQP URL: %s", err)
	}
	slog.Debug(fmt.Sprintf("Connection URL for RabbitMQ is %s", amqpUrl))

	_, err = rabbitmqContainer.State(ctx)
	if err != nil {
		log.Fatalf("failed to get container state: %s", err)
	}

	portNat, _ := rabbitmqContainer.MappedPort(ctx, "5672") //FIXME
	port, _ := strconv.Atoi(portNat.Port())
	slog.Info(fmt.Sprintf("RabbitMQ for Integration Test is running at port %d", port))

	host, _ := rabbitmqContainer.Host(ctx)

	rabbitmqDataSource := messaging.RabbitMQDataSource{
		Host:     host,
		Port:     port,
		User:     rabbitmqContainer.AdminUsername,
		Password: rabbitmqContainer.AdminPassword,
		AmqpURL:  amqpUrl,
	}

	RabbitMQDataSource = rabbitmqDataSource

}
