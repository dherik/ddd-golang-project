package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/dherik/ddd-golang-project/internal/app"
	"github.com/dherik/ddd-golang-project/internal/config"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
	"github.com/dherik/ddd-golang-project/internal/messaging/rabbitmq"
)

func main() {
	// Load application configuration
	cfg, err := config.LoadConfig("config.yaml")
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize database connection
	db, err := persistence.NewPostgresDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Initialize RabbitMQ connection
	rabbitMQ, err := rabbitmq.NewRabbitMQ(cfg.RabbitMQURL)
	if err != nil {
		log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close()

	// Initialize application services
	appService := app.NewService(db, rabbitMQ)

	// Initialize HTTP server and routes
	router := app.SetupHTTPRoutes(appService)
	httpAddr := fmt.Sprintf(":%d", cfg.HTTPPort)

	// Start the HTTP server
	fmt.Printf("Starting HTTP server on port %d...\n", cfg.HTTPPort)
	if err := http.ListenAndServe(httpAddr, router); err != nil {
		log.Fatalf("HTTP server error: %v", err)
	}
}
