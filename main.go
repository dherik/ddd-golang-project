package main

import (
	"os"

	"github.com/dherik/ddd-golang-project/internal/app"
	"github.com/dherik/ddd-golang-project/internal/config"
	"github.com/dherik/ddd-golang-project/internal/infrastructure/persistence"
)

func main() {

	dataSource := persistence.Datasource{
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

	server := app.Server{
		Datasource: dataSource,
		HTTPPort:   cfg.HTTPPort,
		JWTSecret:  cfg.JWTSecret,
	}
	server.Start()
}
