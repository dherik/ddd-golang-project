package setup

func SetupIntegrationTest() {
	SetupDatabase()
	// SetupRabbitMQ()
	StartServer(Datasource) //FIXME
}

func StopIntegrationTest() {
	StopDatabase()
}
