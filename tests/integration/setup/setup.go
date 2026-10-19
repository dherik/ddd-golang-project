package setup

// SetupIntegrationTest brings up the infrastructure and the HTTP server the
// suites exercise. It returns an error when any of that fails, so the suite
// stops before running assertions against a server that is not there.
func SetupIntegrationTest() error {
	SetupDatabase()
	SetupRabbitMQ()
	return StartServer(Datasource, RabbitMQDataSource)
}

func StopIntegrationTest() {
	StopDatabase()
}
