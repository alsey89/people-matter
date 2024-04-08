package main

import (
	"verve-hrms/pkg/configs"
	"verve-hrms/pkg/logger"
	"verve-hrms/pkg/postgres_connector"
	"verve-hrms/pkg/server"

	"go.uber.org/fx"
)

var config *configs.Config

func init() {
	config = configs.NewConfig("server")
}

func main() {
	app := fx.New(
		fx.Supply(config),

		logger.Module(),
		server.Module("server"),
		postgres_connector.Module("database"),
	)

	app.Run()
}
