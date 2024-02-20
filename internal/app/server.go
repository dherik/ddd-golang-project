package app

import (
	"context"
	"log/slog"
	"os"

	"github.com/dherik/ddd-golang-project/internal/app/api"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Server struct {
	Datasource  persistence.Datasource
	ServerReady chan bool
}

func (s *Server) Start() {

	taskRepository := persistence.NewRepository(s.Datasource)
	// taskRepository := persistence.NewMemoryRepository()
	taskService := api.NewTaskService(taskRepository)
	taskHandler := api.NewTaskHandler(*taskService)
	loginHandler := api.NewLoginHandler()
	routes := api.NewRouter(taskHandler, loginHandler)

	echo := echo.New()
	setupSlog(echo)

	routes.SetupRoutes(echo)
	echo.Logger.Fatal(echo.Start(":3333"))
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
