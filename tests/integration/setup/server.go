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
	setupServerOnce sync.Once
)

func StartServer(dataSource persistence.Datasource, rabbitMQDataSource rabbitmq.RabbitMQDataSource) {
	setupServerOnce.Do(func() {
		slog.Info("Initializing the HTTP server...")

		server := app.Server{
			Datasource:         dataSource,
			RabbitMQDataSource: rabbitMQDataSource,
			HTTPPort:           3333,
		}

		go server.Start()
		err := waitServiceStart("http://localhost:3333", 20, 100*time.Millisecond)
		if err != nil {
			slog.Error("failed waiting HTTP server to start", "err", err)
			return
		}
	})
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

		slog.Info("service not healthy yet, retrying", "attempt", i+1, "retryInterval", retryInterval)

		// Wait for the specified interval before retrying
		time.Sleep(retryInterval)
	}

	// Return an error if the service does not become healthy within the specified retries
	return fmt.Errorf("service did not become healthy after %d attempts", maxRetries)
}
