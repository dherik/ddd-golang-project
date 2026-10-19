package setup

import "fmt"

// SetupIntegrationTest brings up the infrastructure and the HTTP server the
// suites exercise. It returns an error when any of that fails, so the suite
// stops before running assertions against a server that is not there.
//
// Every step reports through this one return. The step that failed is named in
// the message, so the suite's own failure line says which part did not come up.
func SetupIntegrationTest() error {
	if err := SetupDatabase(); err != nil {
		return fmt.Errorf("database setup: %w", err)
	}

	if err := SetupRabbitMQ(); err != nil {
		return fmt.Errorf("broker setup: %w", err)
	}

	if err := StartServer(Datasource, RabbitMQDataSource); err != nil {
		return fmt.Errorf("server setup: %w", err)
	}

	return nil
}

// StopIntegrationTest tears the infrastructure back down.
func StopIntegrationTest() error {
	return StopDatabase()
}
