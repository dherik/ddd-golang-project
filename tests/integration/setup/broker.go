package setup

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/dherik/ddd-golang-project/internal/infrastructure/messaging/rabbitmq"
	"github.com/testcontainers/testcontainers-go"
	// Aliased because it shares its package name with the infrastructure
	// messaging package above.
	tcrabbitmq "github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	setupBrokerOnce sync.Once
	setupBrokerErr  error

	RabbitMQDataSource rabbitmq.RabbitMQDataSource

	// brokerContainer and brokerCtx stay package-level because
	// PublishMessage and DeadLetteredPayloads exec into the container from
	// test bodies, long after SetupRabbitMQ has returned.
	brokerContainer *tcrabbitmq.RabbitMQContainer
	brokerCtx       context.Context
)

// SetupRabbitMQ starts the broker container the suites exercise. The container
// starts plain: the application declares its complete topology at startup, so
// nothing is pre-seeded here.
//
// It returns an error rather than exiting the process, so the suite fails with
// its own name attached and the rest of the run still reports. The one-time
// initialisation is preserved, so a failure is memoised and returned identically
// on every later call.
func SetupRabbitMQ() error {
	setupBrokerOnce.Do(func() {
		setupBrokerErr = setupRabbitMQ()
	})

	return setupBrokerErr
}

func setupRabbitMQ() error {
	slog.Info("Initializing the RabbitMQ...")

	ctx := context.Background()

	container, err := tcrabbitmq.Run(ctx,
		"rabbitmq:3-management-alpine",
		tcrabbitmq.WithAdminPassword("guest"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Server startup complete").
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		return fmt.Errorf("failed to start the RabbitMQ container: %w", err)
	}

	// Fetching the URL is also the readiness check: the wait strategy above
	// has already returned, so a failure here means the broker is not usable.
	amqpURL, err := container.AmqpURL(ctx)
	if err != nil {
		return fmt.Errorf("failed to get the RabbitMQ connection URL: %w", err)
	}
	slog.Info("RabbitMQ for integration tests is running")

	brokerContainer = container
	brokerCtx = ctx

	// Only the URL is set. The connection-string builder prefers it over the
	// discrete host, port and credentials, so deriving those here would
	// compute a connection nothing ever reads.
	RabbitMQDataSource = rabbitmq.RabbitMQDataSource{
		AmqpURL: amqpURL,
	}

	return nil
}
