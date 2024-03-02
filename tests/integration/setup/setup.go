package setup

func SetupIntegrationTest() {
	SetupDatabase()
	StartServer(Datasource) //FIXME
}

func StopIntegrationTest() {
	StopDatabase()
}
