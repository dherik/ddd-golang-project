package main

import (
	"os"

	"github.com/dherik/ddd-golang-project/internal/app"
	"github.com/dherik/ddd-golang-project/internal/config"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/messaging"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
)

func main() {

	pgsqlDataSource := persistence.Datasource{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,         //FIXME
		User:     "pguser",     //FIXME
		Password: "pgpassword", //FIXME
		Name:     "dddtasks",   //FIXME
	}

	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		panic(err)
	}

	rabbitmqDataSource := messaging.RabbitMQDataSource{
		Host:     "localhost",
		Port:     5672,
		User:     "guest",
		Password: "guest",
	}

	server := app.Server{
		RabbitMQDataSource: rabbitmqDataSource,
		Datasource:         pgsqlDataSource,
		HTTPPort:           cfg.HTTP.Port,
		JWTSecret:          cfg.JWT.Secret,
	}
	server.Start()
}
