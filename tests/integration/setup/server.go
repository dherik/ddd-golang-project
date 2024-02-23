package setup

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/dherik/ddd-golang-project/internal/app"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
)

func StartServer(dataSource persistence.Datasource) {
	serverReady := make(chan bool)

	server := app.Server{
		Datasource:  dataSource,
		ServerReady: serverReady,
	}

	go server.Start()
	waitServiceStart("http://localhost:3333", 20, 100*time.Millisecond)
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
	return fmt.Errorf("Service did not become healthy after %d attempts", maxRetries)
}
