package app

import (
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
	"github.com/labstack/echo/v4"
)

type Server struct {
	Datasource  persistence.Datasource
	ServerReady chan bool
}

func (s *Server) Start() {

	taskRepository := persistence.NewRepository(s.Datasource)
	// taskRepository := persistence.NewMemoryRepository()
	service := NewTaskService(taskRepository)

	e := echo.New()

	SetupHandler(e, service)
	e.Logger.Fatal(e.Start(":3333"))
}
