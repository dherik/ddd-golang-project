package setup

import (
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/dherik/ddd-golang-project/internal/app"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/messaging/rabbitmq"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
)

var (
	initializedServer     bool
	initializedServerLock sync.Mutex
)

func StartServer(dataSource persistence.Datasource, rabbitMQDataSource rabbitmq.RabbitMQDataSource) {

	initializedServerLock.Lock()
	defer initializedServerLock.Unlock()

	if initializedServer {
		slog.Debug("Server it's already initialized")
		return
	}

	slog.Info("Initializing the HTTP server...")

	server := app.Server{
		Datasource:         dataSource,
		RabbitMQDataSource: rabbitMQDataSource,
	}

	go server.Start()
	err := waitServiceStart("http://localhost:3333", 20, 100*time.Millisecond) //FIXME host and port
	if err != nil {
		slog.Error(fmt.Sprintf("failed waiting HTTP server to start: %s", err.Error()))
		return
	}

	initializedServer = true
}

func waitServiceStart(url string, maxRetries int, retryInterval time.Duration) error {
	for i := 0; i < maxRetries; i++ {
		// Make a GET request to the /api/health endpoint
		response, err := http.Get(url + "/api/health")
		if err == nil && response.StatusCode == http.StatusOK {
			// The service is healthy, proceed with your logic
			slog.Info("Service is healthy!")
			return nil
		}

		// Print an error message (optional)
		slog.Info(fmt.Sprintf("Attempt %d: Service is not healthy yet. Retrying in %s...\n", i+1, retryInterval))

		// Wait for the specified interval before retrying
		time.Sleep(retryInterval)
	}

	// Return an error if the service does not become healthy within the specified retries
	return fmt.Errorf("service did not become healthy after %d attempts", maxRetries)
}
