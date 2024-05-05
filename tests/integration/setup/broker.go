package setup

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/rabbitmq"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	initializedBroker     bool
	initializedBrokerLock sync.Mutex
	// PostgresContainer *postgres.PostgresContainer
	// Datasource        persistence.Datasource
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
		testcontainers.WithImage("rabbitmq:3.7.25-management-alpine"),
		// testcontainers.WithWaitStrategy(wait.ForHealthCheck()),
		testcontainers.WithWaitStrategy(
			wait.ForLog("Server startup complete; 3 plugins started.").
				// WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		log.Fatalf("failed to start RabbitMQ container: %s", err)
	}

	state, err := rabbitmqContainer.State(ctx)
	if err != nil {
		log.Fatalf("failed to get container state: %s", err) // nolint:gocritic
	}

	fmt.Println(state.Running)

	// rabbitmqContainer.Start(ctx)

}
