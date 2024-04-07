package main

import (
	"verve-hrms/pkg/logger"
	"verve-hrms/pkg/server"

	"go.uber.org/fx"
)

func main() {
	app := fx.New(
		logger.Module(),
		server.Module("server"),
	)

	app.Run()
}
