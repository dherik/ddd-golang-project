package main

import (
	"log"
	"os"

	"github.com/dherik/ddd-golang-project/internal/app"
	"github.com/dherik/ddd-golang-project/internal/config"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/messaging/rabbitmq"
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
		log.Fatalf("Failed to load the config: %v", err)
	}

	rabbitmqDataSource := rabbitmq.RabbitMQDataSource{
		Host:     os.Getenv("BROKER_HOST"),
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
