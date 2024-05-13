package setup

func SetupIntegrationTest() {
	SetupDatabase()
	SetupRabbitMQ()
	StartServer(Datasource, RabbitMQDataSource) //FIXME
}

func StopIntegrationTest() {
	StopDatabase()
}
