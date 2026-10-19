package setup

import (
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/dherik/ddd-golang-project/internal/app"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/messaging/rabbitmq"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
)

var (
	setupServerOnce sync.Once

	// serverURL is the base URL of the server started by StartServer. The port
	// is chosen by the kernel when the listener is created, so the suite never
	// collides with a locally running app, the docker-compose stack, or a stale
	// container holding a fixed port.
	serverURL string

	// serverStartErr records why the server failed to come up, so the suite
	// fails at setup instead of quietly testing whatever else answers the port.
	serverStartErr error
)

// ServerURL returns the base URL of the server started by StartServer, for
// example "http://127.0.0.1:41234". It is only valid once StartServer has
// returned without error.
func ServerURL() string {
	return serverURL
}

// StartServer binds an ephemeral port, starts the HTTP server on it, and waits
// until it answers. It returns an error if the port cannot be reserved or the
// server never becomes healthy, so callers can fail fast rather than send
// requests to a stranger's listener.
func StartServer(dataSource persistence.Datasource, rabbitMQDataSource rabbitmq.RabbitMQDataSource) error {
	setupServerOnce.Do(func() {
		slog.Info("Initializing the HTTP server...")

		// Bind before starting so the effective port is known up front.
		listener, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			serverStartErr = fmt.Errorf("failed reserving a local port: %w", err)
			return
		}

		serverURL = "http://" + listener.Addr().String()

		server := app.Server{
			Datasource:         dataSource,
			RabbitMQDataSource: rabbitMQDataSource,
			JWTSecret:          "integration-test-secret",
			Listener:           listener,
		}

		go server.Start()

		if err := waitServiceStart(serverURL, 20, 100*time.Millisecond); err != nil {
			serverStartErr = fmt.Errorf("failed waiting HTTP server to start at %s: %w", serverURL, err)
		}
	})

	return serverStartErr
}

func waitServiceStart(url string, maxRetries int, retryInterval time.Duration) error {
	for i := range maxRetries {
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
