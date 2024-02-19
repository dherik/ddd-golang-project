package main

import (
	"os"

	"github.com/dherik/ddd-golang-project/internal/app"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
)

func main() {

	dataSource := persistence.Datasource{
		Host:     os.Getenv("DB_HOST"),
		Port:     5432,
		User:     "pguser",
		Password: "pgpassword",
		Name:     "dddtasks",
	}
	server := app.Server{Datasource: dataSource}
	server.Start()

	// // Load application configuration
	// cfg, err := config.LoadConfig("config.yaml")

	// slog.Info("Port being used", slog.Int("Port", cfg.HTTPPort))
	// if err != nil {
	// 	slog.Error("Failed to load configuration: %v", err)
	// }

	// // Initialize database connection
	// // db, err := persistence.NewPostgresDB(cfg.DatabaseURL)
	// // if err != nil {
	// // 	log.Fatalf("Failed to initialize database: %v", err)
	// // }
	// // defer db.Close()

	// // Initialize RabbitMQ connection
	// // rabbitMQ, err := infrastructure.NewRabbitMQ(cfg.RabbitMQURL)
	// // if err != nil {
	// // 	log.Fatalf("Failed to initialize RabbitMQ: %v", err)
	// // }
	// // defer rabbitMQ.Close()

	// // Initialize application services
	// // appService := app.NewService(db, rabbitMQ)

	// dataSource := persistence.Datasource{Host: os.Getenv("DB_HOST"), Port: 5432, User: "pguser", Password: "pgpassword", Name: "dddtasks"}
	// taskRepository := persistence.NewRepository(dataSource)
	// // taskRepository := persistence.NewMemoryRepository()
	// service := app.NewTaskService(taskRepository)

	// e := echo.New()

	// app.SetupHandler(e, service)
	// e.Logger.Fatal(e.Start(":3333"))
}
