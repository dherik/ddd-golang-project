package main

import (
	"log"

	"github.com/dherik/ddd-golang-project/internal/app"
	"github.com/dherik/ddd-golang-project/internal/config"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/messaging/rabbitmq"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
)

func main() {
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load the config: %v", err)
	}

	pgsqlDataSource := persistence.Datasource{
		Host:     cfg.Database.Host,
		Port:     cfg.Database.Port,
		User:     cfg.Database.User,
		Password: cfg.Database.Password,
		Name:     cfg.Database.Name,
	}

	rabbitmqDataSource := rabbitmq.RabbitMQDataSource{
		Host:     cfg.Broker.Host,
		Port:     cfg.Broker.Port,
		User:     cfg.Broker.User,
		Password: cfg.Broker.Password,
	}

	server := app.Server{
		RabbitMQDataSource: rabbitmqDataSource,
		Datasource:         pgsqlDataSource,
		HTTPPort:           cfg.HTTP.Port,
		JWTSecret:          cfg.JWT.Secret,
	}

	server.Start()
}
