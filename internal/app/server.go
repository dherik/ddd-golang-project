package app

import (
	"context"
	"log"
	"log/slog"
	"os"

	"github.com/dherik/ddd-golang-project/internal/app/api"
	listener "github.com/dherik/ddd-golang-project/internal/app/messaging"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/messaging/rabbitmq"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Datasource         persistence.Datasource
	RabbitMQDataSource rabbitmq.RabbitMQDataSource
	HTTPPort           int
	JWTSecret          string
}

func (s *Server) Start() {

	pgsql := persistence.NewPostgreRepository(s.Datasource)
	taskRepository := persistence.NewTaskRepository(pgsql)
	userRepository := persistence.NewUserRepository(pgsql)
	taskService := api.NewTaskService(taskRepository)
	userService := api.NewUserService(userRepository)
	taskHandler := api.NewTaskHandler(*taskService)
	loginHandler := api.NewLoginHandler(userService, s.JWTSecret)
	userHandler := api.NewUserHandler(userService)
	routes := api.NewRouter(taskHandler, loginHandler, userHandler, s.JWTSecret)

	rabbitMQ, err := rabbitmq.NewRabbitMQ(s.RabbitMQDataSource.ConnectionString()) //FIXME error
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitMQ.Close() //TODO test it

	// calendarQueue := messaging.NewCalendarQueue(rabbitMQ)
	// calendarQueue.StartListenEvents()

	// Create application context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queueCalendar, _ := rabbitMQ.DeclareQueueAndBind(
		"todo-service.events.calendar.birthday", "calendar", "birthday")

	messageListener := listener.NewMessageListener(rabbitMQ)
	calendarHandler := &listener.CalendarHandler{}

	// Start listening to queues
	if err := messageListener.StartListening(ctx, queueCalendar.Name, calendarHandler); err != nil {
		log.Fatalf("Failed to start listening to order queue: %v", err)
	}

	echo := echo.New()
	setupSlog(echo)

	routes.SetupRoutes(echo)
	echo.Logger.Fatal(echo.Start(":3333")) //FIXME port
}

func setupSlog(e *echo.Echo) {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	slog.SetDefault(logger)

	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogError:    true,
		HandleError: true, // forwards error to the global error handler, so it can decide appropriate status code
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
}
